-- name: CreateWindowType :exec
INSERT INTO window_type (
  name, 
  version
) VALUES (
  sqlc.arg('name'),
  sqlc.arg('version')
) ON CONFLICT (name, version) DO NOTHING;

-- name: CreateAlgorithm :exec
WITH processor_id AS (
  SELECT id FROM processor p
  WHERE p.name = sqlc.arg('processor_name') 
  AND p.runtime = sqlc.arg('processor_runtime')
),
window_type_id AS (
  SELECT id FROM window_type w
  WHERE w.name = sqlc.arg('window_type_name') 
  AND w.version = sqlc.arg('window_type_version')
)
INSERT INTO algorithm (
  name,
  version,
  processor_id,
  window_type_id
) VALUES (
  sqlc.arg('name'),
  sqlc.arg('version'),
  (SELECT id FROM processor_id),
  (SELECT id FROM window_type_id)
) ON CONFLICT (name, version, processor_id) DO UPDATE
SET
  window_type_id = excluded.window_type_id;

-- name: ReadAlgorithmsForWindow :many
SELECT a.* FROM algorithm a
JOIN window_type wt ON a.window_type_id = wt.id
WHERE wt.name = sqlc.arg('window_type_name') 
AND wt.version = sqlc.arg('window_type_version');

-- name: CreateAlgorithmDependency :exec
WITH from_algo AS (
  SELECT a.id, a.window_type_id FROM algorithm a
  JOIN processor p ON a.processor_id = p.id
  WHERE a.name = sqlc.arg('from_algorithm_name')
  AND a.version = sqlc.arg('from_algorithm_version')
  AND p.name = sqlc.arg('from_processor_name')
  AND p.runtime = sqlc.arg('from_processor_runtime')
),
to_algo AS (
  SELECT a.id, a.window_type_id FROM algorithm a
  JOIN processor p ON a.processor_id = p.id
  WHERE a.name = sqlc.arg('to_algorithm_name')
  AND a.version = sqlc.arg('to_algorithm_version')
  AND p.name = sqlc.arg('to_processor_name')
  AND p.runtime = sqlc.arg('to_processor_runtime')
)
INSERT INTO algorithm_dependency (
  from_algorithm_id,
  to_algorithm_id,
  path,
  from_window_type_id,
  to_window_type_id
) VALUES (
  (SELECT id FROM from_algo LIMIT 1),
  (SELECT id FROM to_algo LIMIT 1),
  (SELECT text2ltree(from_algo.id::text || '.' || to_algo.id::text)
    FROM from_algo, to_algo),
  (SELECT window_type_id FROM from_algo LIMIT 1),
  (SELECT window_type_id FROM to_algo LIMIT 1)
) ON CONFLICT (from_algorithm_id, to_algorithm_id) DO UPDATE
  SET
    from_window_type_id = excluded.from_window_type_id,
    to_window_type_id = excluded.to_window_type_id;

-- name: ReadAlgorithmDependencies :many
SELECT ad.* FROM algorithm_dependency ad WHERE ad.from_algorithm_id = sqlc.arg('algorithm_id');

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
  RETURNING id
)
-- clean up old algorithm associations
DELETE FROM processor_algorithm
WHERE processor_id = (
  SELECT id FROM processor p
  WHERE p.name = sqlc.arg('name') 
  AND p.runtime = sqlc.arg('runtime')
);

-- name: AddProcessorAlgorithm :exec
WITH processor_id AS (
  SELECT id FROM processor p
  WHERE p.name = sqlc.arg('processor_name') 
  AND p.runtime = sqlc.arg('processor_runtime')
),
algorithm_id AS (
  SELECT id FROM algorithm a
  WHERE a.name = sqlc.arg('algorithm_name') 
  AND a.version = sqlc.arg('algorithm_version')
)
INSERT INTO processor_algorithm (
  processor_id,
  algorithm_id
) VALUES (
  (SELECT id FROM processor_id),
  (SELECT id FROM algorithm_id)
);

-- name: RegisterWindow :one
WITH window_type_id AS (
  SELECT id FROM window_type 
  WHERE name = sqlc.arg('window_type_name') 
  AND version = sqlc.arg('window_type_version')
)
INSERT INTO windows (
  window_type_id,
  time_from, 
  time_to,
  origin
) VALUES (
  (SELECT id FROM window_type_id),
  sqlc.arg('time_from'),
  sqlc.arg('time_to'),
  sqlc.arg('origin')
) RETURNING *;


-- name: GetDag :many
-- SELECT a.id FROM algorithm a
--   WHERE a.window_type_name=sqlc.arg('window_type_name')
--   AND a.window_type_version=sqlc.arg('window_type_version')
--
-- SELECT ad.to_algorithm_id FROM algorithm_dependency ad
-- WHERE ad.from_algorithm_id a.id


