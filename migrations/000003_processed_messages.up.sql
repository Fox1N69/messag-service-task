CREATE TABLE processed_messages (
    id SERIAL PRIMARY KEY,
    message_id INTEGER REFERENCES messages(id),
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    kafka_topic VARCHAR(255),
    kafka_partition INTEGER,
    kafka_offset BIGINT
);
