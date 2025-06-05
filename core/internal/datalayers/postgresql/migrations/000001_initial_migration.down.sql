-- Drop indexes first
DROP INDEX IF EXISTS idx_results_windows_id;
DROP INDEX IF EXISTS idx_results_window_type_id;
DROP INDEX IF EXISTS idx_results_algorithm_id;
DROP INDEX IF EXISTS idx_dependency_to_algo;
DROP INDEX IF EXISTS idx_dependency_from_algo;
DROP INDEX IF EXISTS idx_algorithm_window_type_id;
DROP INDEX IF EXISTS idx_algorithm_processor_id;

-- Drop triggers
DROP TRIGGER IF EXISTS refresh_algorithm_execution_paths_after_dependency_change ON algorithm_dependency;
DROP TRIGGER IF EXISTS refresh_algorithm_execution_paths_after_algorithm_change ON algorithm;

-- Drop function
DROP FUNCTION IF EXISTS refresh_algorithm_exec_paths();

-- Drop materialized view
DROP MATERIALIZED VIEW IF EXISTS algorithm_execution_paths;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS results;
DROP TABLE IF EXISTS algorithm_required_datagetters;
DROP TABLE IF EXISTS windows;
DROP TABLE IF EXISTS algorithm_dependency;
DROP TABLE IF EXISTS algorithm;
DROP TABLE IF EXISTS data_getters;
DROP TABLE IF EXISTS processor;
DROP TABLE IF EXISTS window_type;

-- Drop extension
DROP EXTENSION IF EXISTS ltree;
