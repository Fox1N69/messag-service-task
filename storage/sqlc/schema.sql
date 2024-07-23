CREATE TABLE IF NOT EXISTS messages (
  id BIGSERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  status_id BIGINT REFERENCES message_status(id) DEFAULT 1,
  processed BOOLEAN DEFAULT FALSE,
  kafka_topic VARCHAR(255),
  kafka_partition INTEGER,
  kafka_offset BIGINT
);

CREATE TABLE IF NOT EXISTS message_status (
  id BIGSERIAL PRIMARY KEY,
  status VARCHAR(255) NOT NULL
);

-- Insert default statuses
INSERT INTO message_status (status) VALUES ('Received');
INSERT INTO message_status (status) VALUES ('Processing');
INSERT INTO message_status (status) VALUES ('Processed');

-- Индекс по статусу сообщений
CREATE INDEX idx_messages_status_id ON messages(status_id);

-- Индекс по дате создания сообщений
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- Составной индекс (если часто фильтруете по `status_id` и `created_at`)
CREATE INDEX idx_messages_status_created_at ON messages(status_id, created_at);

-- Индекс по статусу (если часто используете для поиска)
CREATE INDEX idx_message_status_status ON message_status(status);
