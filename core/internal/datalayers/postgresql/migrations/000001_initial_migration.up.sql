CREATE EXTENSION ltree;

-- Window types that can trigger algorithms
CREATE TABLE window_type (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  version TEXT NOT NULL CHECK (version ~ '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$'),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name, version)
);

-- Processors that can execute algorithms
CREATE TABLE processor (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  runtime TEXT NOT NULL, -- e.g. py3.*, go1.*, etc.
  connection_string TEXT NOT NULL, -- the gRPC string to the client
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name, runtime)
);

-- Store of all the algorithms
CREATE TABLE algorithm (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  version TEXT NOT NULL CHECK (version ~ '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$'),
  processor_id BIGINT NOT NULL,
  window_type_id BIGINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name, version),
  FOREIGN KEY (window_type_id) REFERENCES window_type(id),
  FOREIGN KEY (processor_id) REFERENCES processor(id) ON DELETE CASCADE
);

-- Store of all the dependencies between algorithms
CREATE TABLE algorithm_dependency (
  id BIGSERIAL PRIMARY KEY,
  from_algorithm_id BIGINT NOT NULL,
  to_algorithm_id BIGINT NOT NULL, 
  from_window_type_id BIGINT NOT NULL,
  to_window_type_id BIGINT NOT NULL,
  from_processor_id BIGINT NOT NULL,
  to_processor_id BIGINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (from_algorithm_id, to_algorithm_id),
  FOREIGN KEY (from_algorithm_id) REFERENCES algorithm(id) ON DELETE CASCADE,
  FOREIGN KEY (to_algorithm_id) REFERENCES algorithm(id) ON DELETE CASCADE,
  FOREIGN KEY (from_window_type_id) REFERENCES window_type(id),
  FOREIGN KEY (to_window_type_id) REFERENCES window_type(id),
  FOREIGN KEY (from_processor_id) REFERENCES processor(id) ON DELETE CASCADE,
  FOREIGN KEY (to_processor_id) REFERENCES processor(id) ON DELETE CASCADE,
  -- Prevent self-dependencies
  CHECK (from_algorithm_id != to_algorithm_id)
);

-- What algorithms require data from what data getters
CREATE TABLE algorithm_required_datagetters (
  id BIGSERIAL NOT NULL,
  data_getter_id BIGINT NOT NULL, 
  algorithm_id BIGINT NOT NULL,
  FOREIGN KEY (data_getter_id) REFERENCES data_getters(id) ON DELETE CASCADE,
  FOREIGN KEY (algorithm_id) REFERENCES algorithms(id) ON DELETE CASCADE
);

-- Windows that trigger algorithms
CREATE TABLE windows (
  id BIGSERIAL PRIMARY KEY,
  window_type_id BIGINT NOT NULL,
  time_from BIGINT NOT NULL,
  time_to BIGINT NOT NULL,
  origin TEXT NOT NULL, -- the location that emitted the window
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (window_type_id) REFERENCES window_type(id)
);

CREATE TABLE data_getters (
  id bigserial primary key,
  processor_id bigint not null,
  name text not null,
  window_type_id bigint not null,
  ttl_seconds bigint not null,
  max_size_bytes bigint not null,
  foreign key (window_type_id) references window_type(id),
  foreign key (processor_id) references processor(id) ON DELETE CASCADE,
  UNIQUE(processor_id, name)
);

-- Where the results are stored
CREATE TABLE results (
  id BIGSERIAL PRIMARY KEY,
  windows_id BIGINT,
  window_type_id BIGINT, 
  algorithm_id BIGINT, 
  result_value DOUBLE PRECISION,
  result_array DOUBLE PRECISION[],
  result_json JSONB,
  FOREIGN KEY (windows_id) REFERENCES windows(id),
  FOREIGN KEY (window_type_id) REFERENCES window_type(id),
  FOREIGN KEY (algorithm_id) REFERENCES algorithm(id)
);

-- View constructing the algorithm execution paths for the DAG
CREATE MATERIALIZED VIEW IF NOT EXISTS algorithm_execution_paths AS
WITH RECURSIVE leaf_nodes AS (
  -- leaf nodes
    SELECT
        algorithm_dependency.to_algorithm_id
    FROM
        algorithm_dependency
    EXCEPT
    SELECT
        from_algorithm_id
    FROM
        algorithm_dependency
),
search_tree AS (
    -- root nodes
    SELECT
        a.id AS algo_id,
        0 AS num_dependencies,
        a.id::VARCHAR AS algo_id_path,
        a.processor_id::VARCHAR AS proc_id_path,
        a.window_type_id::VARCHAR as window_type_id_path
    FROM
        algorithm a
    WHERE
        a.id NOT IN (
            SELECT ad.to_algorithm_id
            FROM algorithm_dependency ad
        )

    UNION ALL

    SELECT
        ad.to_algorithm_id AS algo_id,
        st.num_dependencies + 1,
        st.algo_id_path || '.' || ad.to_algorithm_id::VARCHAR,
        st.proc_id_path || '.' || ad.to_processor_id::VARCHAR,
        st.window_type_id_path || '.' || ad.to_window_type_id::VARCHAR
    FROM
        algorithm_dependency ad
    JOIN
        search_tree st ON ad.from_algorithm_id = st.algo_id
),
final_view AS (
    SELECT
        st.algo_id AS final_algo_id,
        st.num_dependencies,
        text2ltree(st.algo_id_path) AS algo_id_path,
        text2ltree(st.window_type_id_path) AS window_type_id_path,
        text2ltree(st.proc_id_path) AS proc_id_path
    FROM search_tree st
    WHERE
        st.algo_id IN (SELECT to_algorithm_id FROM leaf_nodes)
        OR st.num_dependencies = 0 -- no dependencies
    ORDER BY nlevel(text2ltree(st.algo_id_path))
)
SELECT * FROM final_view;

-- function to refresh the materialised view
CREATE OR REPLACE FUNCTION refresh_algorithm_exec_paths()
RETURNS TRIGGER AS $$
BEGIN
    REFRESH MATERIALIZED VIEW algorithm_execution_paths;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;


-- triggers to update the view when source tables update
CREATE TRIGGER refresh_algorithm_execution_paths_after_algorithm_change
AFTER INSERT OR UPDATE OR DELETE ON algorithm
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_algorithm_exec_paths();

CREATE TRIGGER refresh_algorithm_execution_paths_after_dependency_change
AFTER INSERT OR UPDATE OR DELETE ON algorithm_dependency
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_algorithm_exec_paths();

-- create some specific indexes to help with query speed
CREATE INDEX idx_algorithm_processor_id ON algorithm(processor_id);
CREATE INDEX idx_algorithm_window_type_id ON algorithm(window_type_id);
CREATE INDEX idx_dependency_from_algo ON algorithm_dependency(from_algorithm_id);
CREATE INDEX idx_dependency_to_algo ON algorithm_dependency(to_algorithm_id);
CREATE INDEX idx_results_algorithm_id ON results(algorithm_id);
CREATE INDEX idx_results_window_type_id ON results(window_type_id);
CREATE INDEX idx_results_windows_id ON results(windows_id);

