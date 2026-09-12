-- name: GetUserPosts :many
-- Returns all posts for a given user
SELECT * FROM posts
WHERE user_id = @user_id
AND (title ILIKE @search OR content ILIKE @search)
ORDER BY id DESC
LIMIT $1 OFFSET $2
;

-- name: GetUserPostsCount :one
-- Returns the count of posts for a given user
SELECT COUNT(*) FROM posts WHERE user_id = @user_id
AND (title ILIKE @search OR content ILIKE @search)
;

-- name: GetPosts :many
-- Returns all posts for a given user
SELECT * FROM posts
WHERE (title ILIKE @search OR content ILIKE @search)
ORDER BY id DESC
LIMIT $1 OFFSET $2
;


-- name: GetPostsCount :one
-- Returns the count of posts for a given user
SELECT COUNT(*) FROM posts
WHERE (title ILIKE @search OR content ILIKE @search)
;

-- name: CreatePost :exec
-- Creates a new post
INSERT INTO posts (user_id, title, content, images)
VALUES (@user_id, @title, @content, @images)
;
