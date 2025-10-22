-- name: UsersGetAll :many
SELECT * from users;

-- name: UsersCreate :one
INSERT INTO users (username, email)
VALUES (:username, :email)
RETURNING *;