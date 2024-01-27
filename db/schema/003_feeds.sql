-- +goose Up
CREATE TABLE
  feeds (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    -- Relationships
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
  );


-- +goose Down
DROP TABLE
  feeds;