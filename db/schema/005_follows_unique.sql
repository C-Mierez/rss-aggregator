-- +goose Up 
ALTER TABLE
  follows
ADD
  CONSTRAINT follows_feed_id_user_id_unique UNIQUE (feed_id, user_id);


-- +goose Down
ALTER TABLE
  follows DROP CONSTRAINT follows_feed_id_user_id_unique;