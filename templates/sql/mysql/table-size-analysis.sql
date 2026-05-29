-- Table Size Analysis: MySQL
SELECT table_schema, table_name,
       ROUND((data_length + index_length) / 1024 / 1024, 2) AS size_mb,
       table_rows AS estimated_rows
FROM information_schema.tables
WHERE table_schema NOT IN ('information_schema','mysql','performance_schema','sys')
ORDER BY (data_length + index_length) DESC
LIMIT 20;
