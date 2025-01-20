package db

import (
	"context"
	"github.com/SemyonTolkachyov/amb/src/common/schema"
	"github.com/jackc/pgx/v5"
	"log"
)

type PostgresRepository struct {
	db *pgx.Conn
}

func NewPostgres(ctx context.Context, url string) (*PostgresRepository, error) {
	db, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close(ctx context.Context) {
	if err := r.db.Close(ctx); err != nil {
		log.Fatal(err)
	}
}

func (r *PostgresRepository) InsertMessage(ctx context.Context, message schema.Message) error {
	_, err := r.db.Exec(ctx, "INSERT INTO messages(id, body, created_at) VALUES($1, $2, $3)", message.Id, message.Body, message.CreatedAt)
	return err
}

func (r *PostgresRepository) ListMessages(ctx context.Context, skip uint64, take uint64) ([]schema.Message, error) {
	rows, err := r.db.Query(ctx, "SELECT id, body, created_at FROM messages ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse all rows into an array of Messages
	var messages []schema.Message
	for rows.Next() {
		message := schema.Message{}
		if err = rows.Scan(&message.Id, &message.Body, &message.CreatedAt); err == nil {
			messages = append(messages, message)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
