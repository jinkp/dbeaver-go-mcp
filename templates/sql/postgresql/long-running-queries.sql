-- Long Running Queries: PostgreSQL
SELECT
    pid, now() - query_start AS duration,
    usename, state, query
FROM pg_stat_activity
WHERE state != 'idle'
  AND query_start IS NOT NULL
ORDER BY duration DESC
LIMIT 20;
