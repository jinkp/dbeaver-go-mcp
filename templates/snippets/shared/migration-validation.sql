-- Migration Validation Template
-- Run on both source and target to compare schemas
-- Replace <schema_name> with your schema
SELECT table_name, column_name, data_type
FROM information_schema.columns
WHERE table_schema = '<schema_name>'
ORDER BY table_name, ordinal_position;
