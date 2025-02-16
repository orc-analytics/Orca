-- name: CreateAlgorithmType :one
INSERT INTO algorithm_types (
    name,
    version,
    window_type,
    depends_on
) VALUES (
    sqlc.arg('name'),
    sqlc.arg('version'),
    sqlc.arg('window_type'),
    sqlc.arg('depends_on')
) RETURNING *;

-- name: CreateWindowType :one
INSERT INTO window_types (
    name
) VALUES (
    sqlc.arg('name')
) ON CONFLICT (name) DO NOTHING
RETURNING *;
