SELECT t.name AS table_name, i.name AS index_name, i.type_desc
FROM sys.indexes i JOIN sys.tables t ON i.object_id = t.object_id
WHERE i.name IS NOT NULL ORDER BY t.name, i.name;
