-- All metadata fields used by window types
CREATE TABLE metadata_fields (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL
);

-- Bridge table to handle metadata field references
CREATE TABLE metadata_fields_references (
    window_type_id BIGINT REFERENCES window_type(id) ON DELETE CASCADE,
    metadata_fields_id BIGINT REFERENCES metadata_fields(id) ON DELETE CASCADE,
    PRIMARY KEY (window_type_id, metadata_fields_id)
);

-- Index to aid querying
CREATE INDEX idx_metadata_fields_references_id 
ON metadata_fields_references(metadata_fields_id);
