-- Cross-DB Row Count Comparison Template
-- Replace <table_name> with your target table
-- Run on source and target databases separately
SELECT '<table_name>' AS table_name, COUNT(*) AS row_count FROM <table_name>;
