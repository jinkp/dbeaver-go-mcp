-- Compare Row Counts: SQL Server
SELECT
    s.name AS schema_name,
    t.name AS table_name,
    p.rows AS estimated_rows
FROM sys.tables t
JOIN sys.schemas s ON s.schema_id = t.schema_id
JOIN sys.partitions p ON t.object_id = p.object_id AND p.index_id IN (0,1)
ORDER BY p.rows DESC;
