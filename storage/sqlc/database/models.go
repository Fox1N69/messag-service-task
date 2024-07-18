// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Message struct {
	ID        int64
	Content   string
	CreatedAt pgtype.Timestamptz
	StatusID  pgtype.Int8
}

type MessageStatus struct {
	ID     int64
	Status string
}

type ProcessedMessage struct {
	ID             int64
	MessageID      pgtype.Int8
	ProcessedAt    pgtype.Timestamptz
	KafkaTopic     pgtype.Text
	KafkaPartition pgtype.Int4
	KafkaOffset    pgtype.Int8
}
