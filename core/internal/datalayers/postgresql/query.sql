---------------------- Core Operations ----------------------  
-- name: CreateProcessorAndPurgeAlgos :exec
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
  connection_string = EXCLUDED.connection_string
RETURNING id;

-- name: CreateWindowType :exec
INSERT INTO window_type (
  name, 
  version, 
  description
) VALUES (
  sqlc.arg('name'),
  sqlc.arg('version'),
  sqlc.arg('description')
) ON CONFLICT (name, version) DO UPDATE
SET
  name = EXCLUDED.name,
  version = EXCLUDED.version,
  description = EXCLUDED.description;

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
  window_type_id,
  result_type
) VALUES (
  sqlc.arg('name'),
  sqlc.arg('version'),
  (SELECT id FROM processor_id),
  (SELECT id FROM window_type_id),
  sqlc.arg('result_type')
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
  origin, 
  metadata
) VALUES (
  (SELECT id FROM window_type_id),
  sqlc.arg('time_from'),
  sqlc.arg('time_to'),
  sqlc.arg('origin'),
  sqlc.arg('metadata')
) RETURNING window_type_id, id;

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


---------------------- Data operations ---------------------- 
-- name: ReadWindowTypes :many
SELECT
  id, 
  version, 
  name,
  description,
  created
FROM window_type
ORDER BY created DESC;

-- name: ReadAlgorithms :many
SELECT
  a.id,
  a.name,
  a.version,
  a.created,
  a.result_type,
  w.name as window_name, 
  w.version as window_version,
  p.name as processor_name, 
  p.runtime as processor_runtime
FROM algorithm a
  JOIN window_type w ON a.window_type_id = w.id
  JOIN processor p ON a.processor_id = p.id
ORDER BY a.processor_id, a.created DESC;

-- name: ReadProcessors :many
SELECT
  id,
  name, 
  runtime, 
  created
FROM processor
ORDER BY created DESC;

-- name: ReadResultsStats :one
SELECT
  COUNT(r.id)
FROM results r;

-- name: ReadDistinctWindowMetadata :many
SELECT DISTINCT w.metadata FROM windows w
JOIN window_type wt ON w.window_type_id = wt.id 
WHERE
  w.time_from  >= sqlc.arg('time_from')
  AND w.time_to <= sqlc.arg('time_to')
  AND wt.name = sqlc.arg('window_type_name')
  AND wt.version = sqlc.arg('window_type_version');

-- name: ReadResultsForAlgorithm :many
select
  w.time_from,
  w.time_to,
  r.result_value,
  r.result_array,
  r.result_json
from results r
join algorithm a on r.algorithm_id = a.id
join windows w on r.windows_id = w.id
where
	w.time_from  >= sqlc.arg('time_from') and w.time_to <= sqlc.arg('time_to')
	and a."name" = sqlc.arg('algorithm_name')
	and a."version" = sqlc.arg('algorithm_version')
ORDER BY w.time_from, w.time_to ASC;

-- name: ReadDistinctJsonResultFieldsForAlgorithm :many
select distinct jsonb_object_keys(r.result_json) as field_names from results r
join algorithm a on r.algorithm_id = a.id
join windows w on r.windows_id = w.id
where
	w.time_from  >= sqlc.arg('time_from') and w.time_to <= sqlc.arg('time_to')
	and a."name" = sqlc.arg('algorithm_name')
	and a."version" = sqlc.arg('algorithm_version');

-- name: ReadAlgorithmJsonField :many
select w.time_from, w.time_to, (r.result_json::json->>sqlc.arg('field_name')) as result from results r
join algorithm a on r.algorithm_id = a.id
join windows w on r.windows_id = w.id
where
	w.time_from  >= sqlc.arg('time_from') and w.time_to <= sqlc.arg('time_to')
	and a."name" = sqlc.arg('algorithm_name')
	and a."version" = sqlc.arg('algorithm_version');

-- name: ReadWindows :many
select
  w.time_from,
  w.time_to,
  w.origin,
  w.metadata,
  wt.name,
  wt.version
from windows w
join window_type wt on w.window_type_id =wt.id
where
	wt."name" = sqlc.arg('window_type_name') and wt."version" = sqlc.arg('window_type_version')
	and w.time_from  >= sqlc.arg('time_from') and w.time_to <= sqlc.arg('time_to')
ORDER BY w.time_from, w.time_to ASC;

-- name: ReadWindowsForMetadata :many
select
  w.time_from,
  w.time_to,
  w.origin,
  w.metadata,
  wt.name,
  wt.version
from windows w
join window_type wt on w.window_type_id =wt.id
where
	wt."name" = sqlc.arg('window_type_name') and wt."version" = sqlc.arg('window_type_version')
	and w.time_from  >= sqlc.arg('time_from') and w.time_to <= sqlc.arg('time_to')
	and w.metadata::jsonb @> sqlc.arg('metadata_filter')::jsonb
ORDER BY w.time_from, w.time_to ASC;


-- name: ReadResultsForAlgorithmAndMetadata :many
WITH algorithmId AS (
  SELECT id FROM algorithm WHERE name = sqlc.arg('algorithm_name') AND version = sqlc.arg('algorithm_version')
)
SELECT
  w.time_from,
  w.time_to,
  w.metadata,
  r.result_value,
  r.result_array,
  r.result_json
FROM results r
JOIN windows w ON r.windows_id  = w.id
WHERE
  w.time_from >= sqlc.arg('time_from') AND
  w.time_to <= sqlc.arg('time_to') AND
  w.metadata::jsonb @> sqlc.arg('metadata_filter')::jsonb AND
  r.algorithm_id = (SELECT id FROM algorithmId)
ORDER BY w.time_from, w.time_to ASC;



---------------------- Annotation operations ---------------------- 
-- name: CreateAnnotation :exec
WITH new_annotation AS (
  INSERT INTO annotations (time_from, time_to, description) 
  VALUES (sqlc.arg('time_from'), sqlc.arg('time_to'), sqlc.arg('description'))
  RETURNING id
),
algorithm_inserts AS (
  INSERT INTO annotation_algorithms (annotation_id, algorithm_id)
  SELECT na.id, a.id
  FROM new_annotation na
  CROSS JOIN algorithm a
  WHERE a.name IN (sqlc.slice('captured_algorithm_names'))
    AND a.version IN (sqlc.slice('captured_algorithm_versions'))
  RETURNING annotation_id
)
INSERT INTO annotation_window_types (annotation_id, window_type_id)
SELECT na.id, wt.id
FROM new_annotation na
CROSS JOIN window_type wt
WHERE wt.name IN (sqlc.slice('captured_window_names'))
  AND wt.version IN (sqlc.slice('captured_window_versions'));
