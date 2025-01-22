package main

import (
	"context"
	"github.com/SemyonTolkachyov/amb/src/common/db"
	"github.com/SemyonTolkachyov/amb/src/common/event"
	"github.com/SemyonTolkachyov/amb/src/common/schema"
	"github.com/SemyonTolkachyov/amb/src/common/util"
	"github.com/SemyonTolkachyov/amb/src/querysrv/search"
	"log"
	"net/http"
	"strconv"
)

func onMessageCreated(m event.MessageCreatedEvent) {
	// Index message for searching
	message := schema.Message{
		Id:        m.ID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
	if err := search.InsertMessage(context.Background(), message); err != nil {
		log.Println(err)
	}
}

func listMessagesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	// Read parameters
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	// Fetch messages
	messages, err := db.ListMessages(ctx, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch messages")
		return
	}

	util.ResponseOk(w, messages)
}

func searchMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Read parameters
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing query parameter")
		return
	}
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	// Search messages
	messages, err := search.GetMessages(ctx, query, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []schema.Message{})
		return
	}

	util.ResponseOk(w, messages)
}
