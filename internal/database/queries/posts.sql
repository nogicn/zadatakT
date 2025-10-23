-- name: PostsGetAll :many
SELECT * from posts;

-- name: PostsCreate :one
INSERT INTO posts (user_id, title, content)
VALUES (:user_id, :title, :content)
RETURNING *;

-- name: PostsGetByID :one
SELECT * from posts WHERE id = sqlc.arg(id);

-- name: PostsGetByUserID :many
SELECT * from posts WHERE user_id = sqlc.arg(user_id);