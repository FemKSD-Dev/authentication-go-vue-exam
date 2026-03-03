-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  delete_at TIMESTAMPTZ
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_delete_at ON users(delete_at);

-- +goose Down
DROP TABLE IF EXISTS users;
