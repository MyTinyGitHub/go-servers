-- +goose Up

create table chirps(
  id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  body text NOT NULL,
  user_id uuid REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP table chirps;
