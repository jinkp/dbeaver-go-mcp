-- Failed Jobs: SQL Server Agent
SELECT
    j.name AS job_name,
    h.run_date, h.run_time,
    h.run_status,
    h.message
FROM msdb.dbo.sysjobhistory h
JOIN msdb.dbo.sysjobs j ON h.job_id = j.job_id
WHERE h.run_status = 0
ORDER BY h.run_date DESC, h.run_time DESC;
