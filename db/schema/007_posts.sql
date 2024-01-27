-- +goose Up
CREATE TABLE
  posts (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
  );


-- +goose Down
DROP TABLE
  posts;