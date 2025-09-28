-- Drop the new unique constraint
ALTER TABLE algorithm DROP CONSTRAINT algorithm_name_version_window_processor_key;

-- Restore the original unique constraint on (name, version)
ALTER TABLE algorithm ADD CONSTRAINT algorithm_name_version_key 
  UNIQUE (name, version);
