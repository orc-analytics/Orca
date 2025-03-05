-- Window types that can trigger algorithms
CREATE TABLE window_type (
  name TEXT,
  version TEXT NOT NULL CHECK (version ~ '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$'),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (name, version)
);

-- Processors that can execute algorithms
CREATE TABLE processor (
  name TEXT NOT NULL,
  runtime TEXT NOT NULL, -- e.g. py3.*, go1.*, etc.
  connection_string TEXT NOT NULL, -- the gRPC string to the client
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (name, runtime)
);

-- Store of all the algorithms
CREATE TABLE algorithm (
  name TEXT NOT NULL,
  version TEXT NOT NULL CHECK (version ~ '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$'),
  processor_name TEXT NOT NULL,
  processor_runtime TEXT NOT NULL,
  window_type_name TEXT NOT NULL,
  window_type_version TEXT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (name, version, processor_name, processor_runtime),
  FOREIGN KEY (window_type_name, window_type_version) REFERENCES window_type(name, version),
  FOREIGN KEY (processor_name, processor_runtime) REFERENCES processor(name, runtime)
);

-- Store of all the dependencies between algorithms
CREATE TABLE algorithm_dependency (
  from_algorithm_name TEXT NOT NULL,
  from_algorithm_version TEXT NOT NULL,
  from_processor_name TEXT NOT NULL,
  from_processor_runtime TEXT NOT NULL,
  to_algorithm_name TEXT NOT NULL,
  to_algorithm_version TEXT NOT NULL,
  to_processor_name TEXT NOT NULL,
  to_processor_runtime TEXT NOT NULL,
  PRIMARY KEY (from_algorithm_name, from_algorithm_version, from_processor_name, from_processor_runtime, to_algorithm_name, to_algorithm_version, to_processor_name, to_processor_runtime),
  FOREIGN KEY (from_algorithm_name, from_algorithm_version, from_processor_name, from_processor_runtime) REFERENCES algorithm(name, version, processor_name, processor_runtime),
  FOREIGN KEY (to_algorithm_name, to_algorithm_version, to_processor_name, to_processor_runtime) REFERENCES algorithm(name, version, processor_name, processor_runtime),
  -- Prevent self-dependencies
  CHECK (NOT (from_algorithm_name = to_algorithm_name AND from_algorithm_version = to_algorithm_version))
);

-- Map of which processors support which algorithms
CREATE TABLE processor_algorithm (
  processor_name TEXT NOT NULL,
  processor_runtime TEXT NOT NULL,
  algorithm_name TEXT NOT NULL,
  algorithm_version TEXT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
  PRIMARY KEY (processor_name, processor_runtime, algorithm_name, algorithm_version), 
  FOREIGN KEY (algorithm_name, algorithm_version, processor_name, processor_runtime) REFERENCES algorithm(name, version, processor_name, processor_runtime),
  FOREIGN KEY (processor_name, processor_runtime) REFERENCES processor(name, runtime)
);

-- Windows that trigger algorithms
CREATE TABLE windows (
  id BIGSERIAL PRIMARY KEY,
  window_type_name TEXT NOT NULL,
  window_type_version TEXT NOT NULL, 
  time_from BIGINT NOT NULL,
  time_to BIGINT NOT NULL,
  origin TEXT NOT NULL, -- the location that emitted the window
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (window_type_name, window_type_version) REFERENCES window_type(name, version)
);
