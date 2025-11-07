-- name: ListCategories :many
SELECT * FROM categories;

-- name: GetCategory :one
SELECT * FROM categories 
WHERE id = ?;

-- name: CreateCategory :exec
INSERT INTO categories (ID, name, description)
VALUES (?, ?, ?);

-- name: UpdateCategory :exec
UPDATE categories
SET name = ?, description = ?
WHERE id = ?;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ?;

-- name: CreateCourse :exec
INSERT INTO courses (ID, category_id, name, description, price)
VALUES (?, ?, ?, ?, ?);