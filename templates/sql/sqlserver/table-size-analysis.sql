-- Table Size Analysis: SQL Server
SELECT
    s.name AS schema_name,
    t.name AS table_name,
    SUM(a.total_pages) * 8 AS total_kb,
    SUM(a.used_pages) * 8 AS used_kb,
    SUM(a.data_pages) * 8 AS data_kb
FROM sys.tables t
JOIN sys.schemas s ON s.schema_id = t.schema_id
JOIN sys.indexes i ON t.object_id = i.object_id
JOIN sys.partitions p ON i.object_id = p.object_id AND i.index_id = p.index_id
JOIN sys.allocation_units a ON p.partition_id = a.container_id
GROUP BY s.name, t.name
ORDER BY total_kb DESC;
