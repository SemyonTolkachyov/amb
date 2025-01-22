package main

import (
	"github.com/SemyonTolkachyov/amb/src/common/db"
	"github.com/SemyonTolkachyov/amb/src/common/event"
	"github.com/SemyonTolkachyov/amb/src/common/schema"
	"github.com/SemyonTolkachyov/amb/src/common/util"
	"github.com/segmentio/ksuid"
	"html/template"
	"log"
	"net/http"
	"time"
)

func createMessageHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Id string `json:"id"`
	}

	ctx := r.Context()

	// Read parameters
	body := template.HTMLEscapeString(r.FormValue("body"))
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid message length")
		return
	}
	// Create message
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create message")
		return
	}
	message := schema.Message{
		Id:        id.String(),
		Body:      body,
		CreatedAt: createdAt,
	}
	if err := db.InsertMessage(ctx, message); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create message")
		return
	}

	// Publish event
	if err := event.PublishMessageCreated(message); err != nil {
		log.Println(err)
	}

	// Return new message
	util.ResponseOk(w, response{Id: message.Id})
}
