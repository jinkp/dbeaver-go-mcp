SELECT session_id, login_name, status, last_request_start_time, host_name, program_name
FROM sys.dm_exec_sessions WHERE is_user_process = 1 ORDER BY last_request_start_time;
