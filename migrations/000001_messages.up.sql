-- Создание таблицы status
CREATE TABLE IF NOT EXISTS message_status (
  id BIGSERIAL PRIMARY KEY,
  status VARCHAR(255) NOT NULL
);

-- Вставка начальных данных в таблицу status
INSERT INTO message_status (status) VALUES ('Received');
INSERT INTO message_status (status) VALUES ('Processing');
INSERT INTO message_status (status) VALUES ('Processed');

-- Создание таблицы messages
CREATE TABLE IF NOT EXISTS messages (
  id BIGSERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  status_id BIGINT REFERENCES message_status(id),
  processed BOOLEAN DEFAULT FALSE,
  kafka_topic VARCHAR(255),
  kafka_partition INTEGER,
  kafka_offset BIGINT
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_messages_status_id ON messages(status_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_status_created_at ON messages(status_id, created_at);
CREATE INDEX IF NOT EXISTS idx_message_status_status ON message_status(status);

-- Обновление данных для существующих записей
-- Примечание: добавьте сюда SQL-запросы для обновления данных, если необходимо
