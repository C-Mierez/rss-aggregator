-- name: CreateUser :one
INSERT INTO
  users (id, name, created_at, updated_at, api_key)
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    encode(sha256(random():: text:: bytea), 'hex')
  ) RETURNING *;


-- name: GetUserByID :one
SELECT
  *
FROM
  users
WHERE
  id = $1;


-- name: GetUserByAPIKey :one
SELECT
  *
FROM
  users
WHERE
  api_key = $1;


-- name: GetUsers :many
SELECT
  *
FROM
  users
ORDER BY
  id ASC;


-- name: UpdateUser :one
UPDATE
  users
SET
  name = $2,
  updated_at = $3
WHERE
  id = $1 RETURNING *;