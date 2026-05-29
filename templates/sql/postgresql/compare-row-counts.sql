-- Compare Row Counts: PostgreSQL
SELECT
    schemaname,
    tablename,
    n_live_tup AS estimated_row_count
FROM pg_stat_user_tables
ORDER BY n_live_tup DESC;
