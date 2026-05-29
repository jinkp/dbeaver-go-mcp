-- Deadlocks Check: PostgreSQL
SELECT
    pid, usename, application_name, state,
    wait_event_type, wait_event,
    query_start, query
FROM pg_stat_activity
WHERE state != 'idle'
  AND query_start < now() - interval '30 seconds'
ORDER BY query_start;
