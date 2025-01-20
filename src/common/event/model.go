package event

import "time"

type Event interface {
	Key() string
}

type MessageCreatedEvent struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

func (m *MessageCreatedEvent) Key() string {
	return "message.created"
}
