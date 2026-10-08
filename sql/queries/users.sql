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

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    name = $1;
