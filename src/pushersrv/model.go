package main

import "time"

type MessageCreatedModel struct {
	Id        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func newMessageCreatedModel(id string, body string, createdAt time.Time) *MessageCreatedModel {
	return &MessageCreatedModel{
		Id:        id,
		Body:      body,
		CreatedAt: createdAt,
	}
}
