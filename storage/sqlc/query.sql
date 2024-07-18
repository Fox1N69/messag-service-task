-- name: InsertMessage :one
INSERT INTO messages (content, status_id)
VALUES ($1, $2)
RETURNING id;

-- name: GetMessages :many
SELECT id,
  content,
  created_at,
  status_id
FROM messages;

-- name: GetMessageByID :one
SELECT id,
  content,
  created_at,
  status_id
FROM messages
WHERE id = $1;

-- name: UpdateMessageStatus :exec
UPDATE messages
SET status_id = $1
WHERE id = $2;

-- name: InsertProcessedMessage :one
INSERT INTO processed_messages (
    message_id,
    kafka_topic,
    kafka_partition,
    kafka_offset
  )
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetProcessedMessages :many
SELECT pm.id,
  pm.message_id,
  pm.processed_at,
  pm.kafka_topic,
  pm.kafka_partition,
  pm.kafka_offset,
  m.content
FROM processed_messages pm
  JOIN messages m ON pm.message_id = m.id;
