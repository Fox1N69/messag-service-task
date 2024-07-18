CREATE TABLE message_status (
  id SERIAL PRIMARY KEY,
  status VARCHAR(255) NOT NULL
);
INSERT INTO message_status (status)
VALUES ('Received');
INSERT INTO message_status (status)
VALUES ('Processing');
INSERT INTO message_status (status)
VALUES ('Processed');