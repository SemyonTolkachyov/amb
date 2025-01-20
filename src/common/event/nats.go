package event

import (
	"bytes"
	"encoding/gob"
	"github.com/SemyonTolkachyov/amb/src/common/schema"
	"github.com/nats-io/nats.go"
	"log"
)

type NatsEventStore struct {
	nc                         *nats.Conn
	messageCreatedSubscription *nats.Subscription
	messageCreatedChan         chan MessageCreatedEvent
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

func (es *NatsEventStore) SubscribeMessageCreated() (<-chan MessageCreatedEvent, error) {
	m := MessageCreatedEvent{}
	es.messageCreatedChan = make(chan MessageCreatedEvent, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	es.messageCreatedSubscription, err = es.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	// Decode message
	go func() {
		for {
			select {
			case msg := <-ch:
				if err := es.readMessage(msg.Data, &m); err != nil {
					log.Fatal(err)
				}
				es.messageCreatedChan <- m
			}
		}
	}()
	return es.messageCreatedChan, nil
}

func (es *NatsEventStore) OnMessageCreated(f func(MessageCreatedEvent)) (err error) {
	m := MessageCreatedEvent{}
	es.messageCreatedSubscription, err = es.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		if err := es.readMessage(msg.Data, &m); err != nil {
			log.Fatal(err)
		}
		f(m)
	})
	return
}

func (es *NatsEventStore) Close() {
	if es.nc != nil {
		es.nc.Close()
	}
	if es.messageCreatedSubscription != nil {
		if err := es.messageCreatedSubscription.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
	}
	close(es.messageCreatedChan)
}

func (es *NatsEventStore) PublishMessageCreated(message schema.Message) error {
	m := MessageCreatedEvent{message.Id, message.Body, message.CreatedAt}
	data, err := es.encodeEventMessage(&m)
	if err != nil {
		return err
	}
	return es.nc.Publish(m.Key(), data)
}

func (es *NatsEventStore) encodeEventMessage(m Event) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (es *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}
