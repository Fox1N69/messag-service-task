-- name: InsertMessage :one
INSERT INTO messages (content, status_id)
VALUES ($1, $2)
RETURNING id;

-- name: GetMessages :many
SELECT id,
  content,
  created_at,
  status_id,
  processed,
  kafka_topic,
  kafka_partition,
  kafka_offset
FROM messages;

-- name: GetMessageByID :one
SELECT id,
  content,
  created_at,
  status_id,
  processed,
  kafka_topic,
  kafka_partition,
  kafka_offset
FROM messages
WHERE id = $1;

-- name: GetTotalMessages :one
SELECT COUNT(*) AS count FROM messages;

-- name: UpdateMessageStatus :exec
UPDATE messages
SET status_id = $1
WHERE id = $2;

-- name: GetProcessedMessages :many
SELECT id,
  content,
  created_at,
  status_id,
  processed,
  kafka_topic,
  kafka_partition,
  kafka_offset
FROM messages
WHERE processed = TRUE;

-- name: GetProcessedMessagesCount :one
SELECT COUNT(*) AS count
FROM messages
WHERE processed = TRUE;

-- name: UpdateMessageProcessingDetails :exec
UPDATE messages
SET processed = TRUE,
    kafka_topic = $2,
    kafka_partition = $3,
    kafka_offset = $4
WHERE id = $1;


-- name: GetMessageStatistics :one
SELECT
  COUNT(CASE WHEN processed = TRUE THEN 1 END) AS processed_count,
  COUNT(CASE WHEN processed = FALSE AND status_id = 1 THEN 1 END) AS received_count, -- assuming status ID 1 is "Received"
  COUNT(CASE WHEN processed = FALSE AND status_id = 2 THEN 1 END) AS processing_count  -- assuming status ID 2 is "Processing"
FROM messages;