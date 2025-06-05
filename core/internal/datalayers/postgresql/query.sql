
-- TODO Potential exploit to delete processors unintentionally.
-- Would need to factor in some outside source like Git to better
-- manage this.
-- name: DeleteProcessor :exec
DELETE FROM processor WHERE (
  name = sqlc.arg('name') AND
  runtime = sqlc.arg('runtime')
);

-- name: CreateProcessor :exec
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
  runtime = EXCLUDED.runtime,
  connection_string = EXCLUDED.connection_string;

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
) ON CONFLICT DO NOTHING;

-- name: ReadAlgorithmsForWindow :many
SELECT a.* FROM algorithm a
JOIN window_type wt ON a.window_type_id = wt.id
WHERE wt.name = sqlc.arg('window_type_name') 
AND wt.version = sqlc.arg('window_type_version');

-- name: CreateAlgorithmDependency :exec
WITH from_algo AS (
  SELECT a.id, a.window_type_id, a.processor_id FROM algorithm a
  JOIN processor p ON a.processor_id = p.id
  WHERE a.name = sqlc.arg('from_algorithm_name')
  AND a.version = sqlc.arg('from_algorithm_version')
  AND p.name = sqlc.arg('from_processor_name')
  AND p.runtime = sqlc.arg('from_processor_runtime')
),
to_algo AS (
  SELECT a.id, a.window_type_id, a.processor_id FROM algorithm a
  JOIN processor p ON a.processor_id = p.id
  WHERE a.name = sqlc.arg('to_algorithm_name')
  AND a.version = sqlc.arg('to_algorithm_version')
  AND p.name = sqlc.arg('to_processor_name')
  AND p.runtime = sqlc.arg('to_processor_runtime')
)
INSERT INTO algorithm_dependency (
  from_algorithm_id,
  to_algorithm_id,
  from_window_type_id,
  to_window_type_id,
  from_processor_id,
  to_processor_id
) VALUES (
  (SELECT id FROM from_algo LIMIT 1),
  (SELECT id FROM to_algo LIMIT 1),
  (SELECT window_type_id FROM from_algo LIMIT 1),
  (SELECT window_type_id FROM to_algo LIMIT 1),
  (SELECT processor_id FROM from_algo LIMIT 1),
  (SELECT processor_id FROM to_algo LIMIT 1)
) ON CONFLICT (from_algorithm_id, to_algorithm_id) DO UPDATE
  SET
    from_window_type_id = excluded.from_window_type_id,
    to_window_type_id = excluded.to_window_type_id,
    from_processor_id = excluded.from_processor_id,
    to_processor_id = excluded.to_processor_id;

-- name: CreateAlgorithmRequiredDataGetter :exec
WITH algo AS (
  SELECT a.id FROM algorithm a
  WHERE a.name = sqlc.arg('algorithm_name') AND a.version = sqlc.arg('algorithm_version')
),
processor AS (
  SELECT p.id FROM processor p
  WHERE p.name = sqlc.arg('processor_name') AND p.runtime = sqlc.arg('processor_runtime')
),
datagetter AS (
  SELECT d.id FROM data_getters d
  WHERE d.name = sqlc.arg('datagetter_name') AND d.processor_id = (SELECT id FROM processor)
)

INSERT INTO algorithm_required_datagetters (
  data_getter_id,
  algorithm_id
) VALUES (
  (SELECT id FROM datagetter),
  (SELECT id FROM algo)
);

-- name: ReadFromAlgorithmDependencies :many
WITH from_algo AS (
  SELECT a.id, a.window_type_id, a.processor_id FROM algorithm a
  JOIN processor p ON a.processor_id = p.id
  WHERE a.name = sqlc.arg('from_algorithm_name')
  AND a.version = sqlc.arg('from_algorithm_version')
  AND p.name = sqlc.arg('from_processor_name')
  AND p.runtime = sqlc.arg('from_processor_runtime')
)
SELECT ad.* FROM algorithm_dependency ad WHERE ad.from_algorithm_id = from_algo.id;

-- name: ReadAlgorithmId :one
WITH processor_id AS (
  SELECT p.id FROM processor p
  WHERE p.name = sqlc.arg('processor_name')
  AND p.runtime = sqlc.arg('processor_runtime')
)
SELECT a.id FROM algorithm a
WHERE a.name = sqlc.arg('algorithm_name')
AND a.version = sqlc.arg('algorithm_version')
AND a.processor_id = (SELECT id from processor_id);
  
-- name: ReadAlgorithmExecutionPaths :many
SELECT aep.* FROM algorithm_execution_paths aep WHERE aep.window_type_id_path ~ ('*.' || sqlc.arg('window_type_id')::TEXT || '.*')::lquery;

-- name: ReadAlgorithmExecutionPathsForAlgo :many
SELECT aep.* FROM algorithm_execution_paths aep WHERE aep.final_algo_id=sqlc.arg('algo_id');

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
) RETURNING window_type_id, id;

-- name: CreateDataGetter :one
WITH processor_id AS (
  SELECT id FROM processor p
  WHERE p.name = sqlc.arg('processor_name') 
  AND p.runtime = sqlc.arg('processor_runtime')
), window_type_id AS (
  SELECT id FROM window_type w
  WHERE w.name = sqlc.arg('window_name') 
  AND w.version = sqlc.arg('window_version')
)
INSERT INTO data_getters (
  processor_id,
  window_type_id,
  name,
  ttl_seconds,
  max_size_bytes
) VALUES (
  (SELECT id FROM processor_id),
  (SELECT id FROM window_type_id),
  sqlc.arg('name'),
  sqlc.arg('ttl_seconds'),
  sqlc.arg('max_size_bytes')
) ON CONFLICT (name, processor_id) DO UPDATE
SET 
  name = EXCLUDED.name,
  processor_id = EXCLUDED.processor_id
RETURNING id;

-- name: CreateResult :one
INSERT INTO results (
  windows_id,
  window_type_id, 
  algorithm_id, 
  result_value,
  result_array,
  result_json
) VALUES (
  sqlc.arg('windows_id'),
  sqlc.arg('window_type_id'),
  sqlc.arg('algorithm_id'),
  sqlc.arg('result_value'),
  sqlc.arg('result_array'),
  sqlc.arg('result_json')
) RETURNING id;

-- name: ReadAllProcessors :many
SELECT 
  id,
  name,
  runtime,
  connection_string,
  created
FROM processor
ORDER BY name, runtime;

-- name: ReadProcessorsByIDs :many
SELECT 
  id,
  name,
  runtime,
  connection_string,
  created
FROM processor
WHERE id = ANY(sqlc.arg('processor_ids')::bigint[])
ORDER BY name, runtime;
