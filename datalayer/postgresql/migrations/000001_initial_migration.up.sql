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
  UNIQUE (name, version, processor_id),
  FOREIGN KEY (window_type_id) REFERENCES window_type(id),
  FOREIGN KEY (processor_id) REFERENCES processor(id)
);

-- Store of all the dependencies between algorithms
CREATE TABLE algorithm_dependency (
  id BIGSERIAL PRIMARY KEY,
  path ltree NOT NULL, -- the dependency path of the algorithm ids, e.g. 5.6
  from_algorithm_id BIGINT NOT NULL,
  to_algorithm_id BIGINT NOT NULL,
  from_window_type_id BIGINT NOT NULL,
  to_window_type_id BIGINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (from_algorithm_id, to_algorithm_id),
  FOREIGN KEY (from_algorithm_id) REFERENCES algorithm(id),
  FOREIGN KEY (to_algorithm_id) REFERENCES algorithm(id),
  FOREIGN KEY (from_window_type_id) REFERENCES window_type(id),
  FOREIGN KEY (to_window_type_id) REFERENCES window_type(id),
  -- Prevent self-dependencies
  CHECK (from_algorithm_id != to_algorithm_id)
);

-- Map of which processors support which algorithms
CREATE TABLE processor_algorithm (
  id BIGSERIAL PRIMARY KEY,
  processor_id BIGINT NOT NULL,
  algorithm_id BIGINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (processor_id, algorithm_id),
  FOREIGN KEY (algorithm_id) REFERENCES algorithm(id),
  FOREIGN KEY (processor_id) REFERENCES processor(id)
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
