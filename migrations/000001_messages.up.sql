-- Создание таблицы message_status
CREATE TABLE IF NOT EXISTS message_status (
  id BIGSERIAL PRIMARY KEY,
  status VARCHAR(255) NOT NULL
);

-- Вставка стандартных статусов
INSERT INTO message_status (status) VALUES ('Received');
INSERT INTO message_status (status) VALUES ('Processing');
INSERT INTO message_status (status) VALUES ('Processed');

-- Создание таблицы messages
CREATE TABLE IF NOT EXISTS messages (
  id BIGSERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  status_id BIGINT REFERENCES message_status(id)
);

-- Индекс по статусу сообщений
CREATE INDEX IF NOT EXISTS idx_messages_status_id ON messages(status_id);

-- Индекс по дате создания сообщений
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- Составной индекс (если часто фильтруете по `status_id` и `created_at`)
CREATE INDEX IF NOT EXISTS idx_messages_status_created_at ON messages(status_id, created_at);

-- Создание таблицы processed_messages
CREATE TABLE IF NOT EXISTS processed_messages (
  id BIGSERIAL PRIMARY KEY,
  message_id BIGINT REFERENCES messages(id),
  processed_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  kafka_topic VARCHAR(255),
  kafka_partition INTEGER,
  kafka_offset BIGINT
);
