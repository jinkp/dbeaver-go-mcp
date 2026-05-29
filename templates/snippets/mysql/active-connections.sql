SELECT id, user, host, db, command, time, state, info FROM information_schema.processlist
WHERE command != 'Sleep' ORDER BY time DESC;
