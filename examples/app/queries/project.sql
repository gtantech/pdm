-- name: FindAllProjects :many
SELECT * FROM projects;

-- name: FindProjectById :one
SELECT * FROM projects WHERE id = ?;

-- name: FindProjectByName :many
SELECT * FROM projects WHERE disp_name LIKE ?;

-- name: InsertProject :exec
INSERT INTO projects (id, disp_name) VALUES (?, ?);

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = ?;