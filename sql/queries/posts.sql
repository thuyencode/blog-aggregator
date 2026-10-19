-- name: CreatePost :one
INSERT INTO
    posts (
        created_at,
        updated_at,
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES
    (
        pg_catalog.now(),
        pg_catalog.now(),
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING
    *;

-- name: GetPosts :many
SELECT
    *
FROM
    posts
ORDER BY
    published_at DESC
LIMIT
    $1;
