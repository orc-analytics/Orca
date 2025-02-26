-- name: AddProcessor :one
INSERT INTO processors (
  name,
  runtime_id,
  active
) VALUES (
  sqlc.arg('name'),
  sqlc.arg('runtime_id'),
  true
) ON CONFLICT (name) DO UPDATE 
SET 
  runtime_id = EXCLUDED.runtime_id,
  active = EXCLUDED.active
RETURNING *;

-- -- name: CreateAlgorithmType :one
-- INSERT INTO algorithm_types (
--     name,
--     version,
--     window_type_name
-- ) VALUES (
--     sqlc.arg('name'),
--     sqlc.arg('version'),
--     sqlc.arg('window_type_name')
-- ) RETURNING *;

-- -- name: CreateWindowType :one
-- INSERT INTO window_types (
--     name
-- ) VALUES (
--     sqlc.arg('name')
-- ) ON CONFLICT (name) DO NOTHING
-- RETURNING *;
--
-- -- name: CreateAlgorithmDependency :one
-- INSERT INTO algorithm_dependencies (
--     algorithm_name,
--     algorithm_version,
--     depends_on_name,
--     depends_on_version
-- ) VALUES (
--     sqlc.arg('algorithm_name'),
--     sqlc.arg('algorithm_version'),
--     sqlc.arg('depends_on_name'),
--     sqlc.arg('depends_on_version')
-- ) RETURNING *;
