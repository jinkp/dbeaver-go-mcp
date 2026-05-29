-- Failed Jobs: PostgreSQL (pg_cron extension)
-- Requires pg_cron extension
SELECT jobid, jobname, start_time, end_time, status, return_message
FROM cron.job_run_details
WHERE status = 'failed'
ORDER BY start_time DESC
LIMIT 50;
