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
    window_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (window_type) REFERENCES window_types(name)
);

-- Different types of algorithms
CREATE TABLE algorithm_types (
    name TEXT NOT NULL,
    version TEXT NOT NULL CHECK (version ~ '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$'),
    window_type TEXT NOT NULL,
    depends_on TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (name, version),
    FOREIGN KEY (window_type) REFERENCES window_types(name)
);
