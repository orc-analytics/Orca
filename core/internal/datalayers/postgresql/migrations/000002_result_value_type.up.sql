CREATE TYPE result_type AS ENUM ('struct', 'array', 'value', 'none');
ALTER TABLE algorithm ADD IF NOT EXISTS result_type result_type;
