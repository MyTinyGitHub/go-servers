-- +goose Up
CREATE TABLE users (
  id UUID,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  email TEXT,

  PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE users;
