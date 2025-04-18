
-- +goose Up
CREATE TABLE refresh_token (
  token TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id uuid NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP,

  PRIMARY KEY(token),
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_token;
