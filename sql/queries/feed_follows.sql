-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO
        feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES
        (
            pg_catalog.now(),
            pg_catalog.now(),
            $1,
            $2
        )
    RETURNING
        *
)
SELECT
    inserted_feed_follow.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM
    inserted_feed_follow
    JOIN users ON inserted_feed_follow.user_id = users.id
    JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollows :many
SELECT
    feeds.name AS feed_name
FROM
    feed_follows
    JOIN users ON feed_follows.user_id = users.id
    JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    feed_follows.user_id = $1;

-- name: RemoveFeedFollow :one
WITH deleted_feed_follow AS (
    DELETE FROM
        feed_follows
    WHERE
        feed_follows.user_id = $1
        AND feed_follows.feed_id = $2
    RETURNING
        *
)
SELECT
    feeds.name AS feed_name
FROM
    deleted_feed_follow
    JOIN feeds ON deleted_feed_follow.feed_id = feeds.id;
