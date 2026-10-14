-- name: CreateFeed :one
INSERT INTO
    feeds (id, created_at, updated_at, name, url, user_id)
VALUES
    (
        pg_catalog.gen_random_uuid(),
        pg_catalog.now(),
        pg_catalog.now(),
        $1,
        $2,
        $3
    )
RETURNING
    *;
