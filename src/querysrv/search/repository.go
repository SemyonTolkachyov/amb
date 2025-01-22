package search

import (
	"context"
	"github.com/SemyonTolkachyov/amb/src/common/schema"
)

type Repository interface {
	Close()
	InsertMessage(ctx context.Context, message schema.Message) error
	SearchMessages(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Message, error)
}

var impl Repository

func SetRepository(repository Repository) {
	println("SetRepository")
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertMessage(ctx context.Context, message schema.Message) error {
	return impl.InsertMessage(ctx, message)
}

func GetMessages(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Message, error) {
	return impl.SearchMessages(ctx, query, skip, take)
}
