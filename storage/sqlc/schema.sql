-- Messages table
CREATE TABLE IF NOT EXISTS messages (
  id BIGSERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  status_id BIGINT REFERENCES message_status(id)
);

CREATE TABLE IF NOT EXISTS message_status (
  id BIGSERIAL PRIMARY KEY,
  status VARCHAR(255) NOT NULL
);

-- Insert default statuses
INSERT INTO message_status (status) VALUES ('Received');
INSERT INTO message_status (status) VALUES ('Processing');
INSERT INTO message_status (status) VALUES ('Processed');

CREATE TABLE IF NOT EXISTS processed_messages (
  id BIGSERIAL PRIMARY KEY,
  message_id BIGINT REFERENCES messages(id),
  processed_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  kafka_topic VARCHAR(255),
  kafka_partition INTEGER,
  kafka_offset BIGINT
);
