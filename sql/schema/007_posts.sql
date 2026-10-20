-- +goose Up
ALTER TABLE
    posts DROP CONSTRAINT posts_feed_id_fkey;

ALTER TABLE
    posts
ADD
    CONSTRAINT posts_feed_id_fkey FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE
    posts DROP CONSTRAINT posts_feed_id_fkey;

ALTER TABLE
    posts
ADD
    CONSTRAINT posts_feed_id_fkey FOREIGN KEY (feed_id) REFERENCES feeds(id);
