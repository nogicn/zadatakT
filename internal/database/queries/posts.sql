-- name: PostsGetAll :many
SELECT * from posts;

-- name: PostsCreate :one
INSERT INTO posts (user_id, title, content)
VALUES (:user_id, :title, :content)
RETURNING *;

-- name: PostsGetByID :one
SELECT * FROM posts WHERE id = sqlc.arg(id);

-- name: PostsGetByUserID :many
SELECT * FROM posts WHERE user_id = sqlc.arg(user_id);