-- name: CreateUser :one
INSERT INTO
    users (id, created_at, updated_at, name)
VALUES
    (
        pg_catalog.gen_random_uuid(),
        pg_catalog.now(),
        pg_catalog.now(),
        $1
    )
RETURNING
    *;

-- name: GetUserByName :one
SELECT
    *
FROM
    users
WHERE
    name = $1;

-- name: GetUsers :many
SELECT
    *
FROM
    users;

-- name: DeleteAllUsers :exec
DELETE FROM
    users;

-- name: MarkFeedFetched :exec
UPDATE
    feeds
SET
    updated_at = pg_catalog.now(),
    last_fetched_at = pg_catalog.now()
WHERE
    id = $1;

-- name: GetNextFeedToFetch :one
SELECT
    *
FROM
    feeds
ORDER BY
    last_fetched_at ASC NULLS FIRST
LIMIT
    1;
