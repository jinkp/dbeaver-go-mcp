-- Health Check: SQL Server
SELECT
    @@VERSION AS sql_version,
    DB_NAME() AS database_name,
    SYSTEM_USER AS connected_user,
    GETDATE() AS server_time,
    SUM(size * 8 / 1024) AS db_size_mb
FROM sys.database_files;
