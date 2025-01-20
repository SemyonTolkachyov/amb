package event

import "github.com/SemyonTolkachyov/amb/src/common/schema"

type Store interface {
	Close()
	PublishMessageCreated(message schema.Message) error
	SubscribeMessageCreated() (<-chan MessageCreatedEvent, error)
	OnMessageCreated(f func(MessageCreatedEvent)) error
}

var impl Store

func SetEventStore(es Store) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishMessageCreated(message schema.Message) error {
	return impl.PublishMessageCreated(message)
}

func SubscribeMessageCreated() (<-chan MessageCreatedEvent, error) {
	return impl.SubscribeMessageCreated()
}

func OnMessageCreated(f func(MessageCreatedEvent)) error {
	return impl.OnMessageCreated(f)
}
