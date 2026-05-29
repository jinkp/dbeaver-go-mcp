-- Deadlocks Check: SQL Server
SELECT
    r.session_id, r.blocking_session_id,
    r.wait_type, r.wait_time,
    r.status, t.text AS query_text
FROM sys.dm_exec_requests r
CROSS APPLY sys.dm_exec_sql_text(r.sql_handle) t
WHERE r.blocking_session_id > 0;
