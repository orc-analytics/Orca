-- name: CreateWindowType :exec
INSERT INTO window_type (
  name, 
  version
) VALUES (
  sqlc.arg('name'),
  sqlc.arg('version')
) ON CONFLICT (name, version) DO NOTHING;

-- name: CreateAlgorithm :exec
INSERT INTO algorithm (
  name,
  version,
  processor_name,
  processor_runtime,
  window_type_name,
  window_type_version
) VALUES (
  sqlc.arg('name'),
  sqlc.arg('version'),
  sqlc.arg('processor_name'),
  sqlc.arg('processor_runtime'),
  sqlc.arg('window_type_name'),
  sqlc.arg('window_type_version')
) ON CONFLICT (name, version, processor_name, processor_runtime) DO NOTHING;

-- name: ReadAlgorithmsForWindow :many
SELECT * FROM algorithm
WHERE
  window_type_name = sqlc.arg('window_type_name') 
  AND window_type_version = sqlc.arg('window_type_version');
  
-- name: CreateAlgorithmDependency :exec
INSERT INTO algorithm_dependency (
  from_algorithm_name,
  from_algorithm_version,
  from_processor_name,
  from_processor_runtime, 
  to_algorithm_name,
  to_algorithm_version,
  to_processor_name,
  to_processor_runtime
) VALUES (
  sqlc.arg('from_algorithm_name'),
  sqlc.arg('from_algorithm_version'),
  sqlc.arg('from_processor_name'),
  sqlc.arg('from_processor_runtime'),
  sqlc.arg('to_algorithm_name'),
  sqlc.arg('to_algorithm_version'),
  sqlc.arg('to_processor_name'),
  sqlc.arg('to_processor_runtime')
) ON CONFLICT DO NOTHING;

-- name: ReadAlgorithmDependencies :many
SELECT * FROM algorithm_dependency WHERE 
  from_algorithm_name = sqlc.arg('algorithm_name')
  AND from_algorithm_version = sqlc.arg('algorithm_version')
  AND from_processor_name = sqlc.arg('processor_name')
  AND from_processor_runtime = sqlc.arg('processor_runtime');

-- name: CreateProcessorAndPurgeAlgos :exec
WITH processor_insert AS (
  INSERT INTO processor (
    name,
    runtime,
    connection_string
  ) VALUES (
    sqlc.arg('name'),
    sqlc.arg('runtime'),
    sqlc.arg('connection_string')
  ) ON CONFLICT (name, runtime) DO UPDATE 
  SET 
    name = EXCLUDED.name,
    runtime = EXCLUDED.runtime
)

-- clean up old algorithm associations
DELETE FROM processor_algorithm
WHERE processor_name = sqlc.arg('name') AND processor_runtime = sqlc.arg('runtime');

-- name: AddProcessorAlgorithm :exec
INSERT INTO processor_algorithm (
  processor_name,
  processor_runtime,
  algorithm_name,
  algorithm_version
) VALUES (
  sqlc.arg('processor_name'),
  sqlc.arg('processor_runtime'),
  sqlc.arg('algorithm_name'),
  sqlc.arg('algorithm_version')
);


-- name: RegisterWindow :one
INSERT INTO windows (
  window_type_name, 
  window_type_version,
  time_from, 
  time_to,
  origin
) VALUES (
  sqlc.arg('window_type_name'),
  sqlc.arg('window_type_version'),
  sqlc.arg('time_from'),
  sqlc.arg('time_to'),
  sqlc.arg('origin')
) RETURNING *;
