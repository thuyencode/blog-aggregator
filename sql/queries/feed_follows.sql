-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO
        feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES
        (
            pg_catalog.now(),
            pg_catalog.now(),
            (
                SELECT
                    users.id
                FROM
                    users
                WHERE
                    users.name = sqlc.arg(user_name)
            ),
            $1
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

-- name: GetFeedFollowsForUser :many
SELECT
    feeds.name AS feed_name
FROM
    feed_follows
    JOIN users ON feed_follows.user_id = users.id
    JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    feed_follows.user_id = (
        SELECT
            users.id
        FROM
            users
        WHERE
            users.name = sqlc.arg(user_name)
    );
