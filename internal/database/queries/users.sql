-- name: UsersGetAll :many
SELECT * from users;

-- name: UsersCreate :one
INSERT INTO users (username, email)
VALUES (:username, :email)
RETURNING *;

-- name: UsersGetByID :one 
SELECT * from users WHERE id = sqlc.arg(id);

-- name: UsersGetByUsername :one
SELECT * from users WHERE username = sqlc.arg(username);

-- name: UsersGetByEmail :one
SELECT * from users WHERE email = sqlc.arg(email);

-- name: UsersUpdateEmailByID :one
UPDATE users
SET email = :email
WHERE id = :id
RETURNING *;
