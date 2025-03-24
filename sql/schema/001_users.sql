-- +goose Up
CREATE TABLE users (
  id UUID,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  email TEXT NOT NULL,
  hashed_password TEXT NOT NULL,
  is_chirpy_red BOOLEAN NOT NULL DEFAULT false,

  PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE users;
