package schema

import "time"

type Message struct {
	Id        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
