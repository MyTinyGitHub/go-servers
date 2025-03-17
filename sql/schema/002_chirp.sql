-- +goose Up

create table chirps(
  id uuid,
  created_at timestamp not null,
  updated_at timestamp not null,
  body text not null,
  user_id uuid REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP table chirps;
