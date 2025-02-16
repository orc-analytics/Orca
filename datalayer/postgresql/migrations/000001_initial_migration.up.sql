-- Types of windows
CREATE TABLE window_types (
    name TEXT PRIMARY KEY NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Windows that trigger algorithms
CREATE TABLE windows (
    id SERIAL PRIMARY KEY,
    time_from BIGINT NOT NULL,
    time_to BIGINT NOT NULL,
    window_type_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (window_type) REFERENCES window_types(name)
);

-- Store of all the algorithms
CREATE TABLE algorithm_types (
    name TEXT NOT NULL,
    version TEXT NOT NULL CHECK (version ~ '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$'),
    window_type_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (name, version),
    FOREIGN KEY (window_type_name) REFERENCES window_types(name)
);

-- Store of all the dependencies between algorithms
CREATE TABLE algorithm_dependencies (
    algorithm_name TEXT NOT NULL,
    algorithm_version TEXT NOT NULL,
    depends_on_name TEXT NOT NULL,
    depends_on_version TEXT NOT NULL,
    PRIMARY KEY (algorithm_name, algorithm_version, depends_on_name, depends_on_version),
    FOREIGN KEY (algorithm_name, algorithm_version) REFERENCES algorithm_types(name, version),
    FOREIGN KEY (depends_on_name, depends_on_version) REFERENCES algorithm_types(name, version),
    CHECK (NOT (algorithm_name = depends_on_name AND algorithm_version = depends_on_version))
);
