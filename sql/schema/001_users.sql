-- +goose Up
CREATE TABLE users (
  id UUID,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  email TEXT NOT NULL,
  hashed_password TEXT NOT NULL,

  PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE users;
