-- Long Running Queries: SQL Server
SELECT TOP 20
    r.session_id, r.status, r.start_time,
    DATEDIFF(SECOND, r.start_time, GETDATE()) AS duration_sec,
    t.text AS query_text
FROM sys.dm_exec_requests r
CROSS APPLY sys.dm_exec_sql_text(r.sql_handle) t
ORDER BY duration_sec DESC;
