-- sql/schema/002_chirps.sql
-- +goose Up
CREATE TABLE chirps (
  id UUID PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  body TEXT NOT NULL,
  user_id UUID NOT NULL,
  CONSTRAINT fk_chirps_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
