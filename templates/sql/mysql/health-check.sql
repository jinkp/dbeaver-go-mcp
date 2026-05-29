-- Health Check: MySQL
SELECT VERSION() AS mysql_version, DATABASE() AS current_db,
       USER() AS current_user, NOW() AS server_time;
