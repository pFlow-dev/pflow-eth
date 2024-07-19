-- DROP MATERIALIZED VIEW transaction_logs_view;
CREATE MATERIALIZED VIEW transaction_logs_view AS
WITH expanded_logs AS (SELECT transaction_hash,
                              block_number,
                              jsonb_array_elements(logs)                 AS log_entry,
                              transaction_details -> 'result' ->> 'from' AS from_address
                       FROM transactions
                       WHERE logs IS NOT NULL)
SELECT uint256_hex_to_numeric(log_entry ->> 'data')        as sequence,
       block_number,
       uint256_hex_to_numeric(log_entry -> 'topics' ->> 1) as role,
       uint256_hex_to_numeric(log_entry -> 'topics' ->> 2) as action,
       uint256_hex_to_numeric(log_entry -> 'topics' ->> 3) as scalar,
       from_address,
       transaction_hash,
       log_entry -> 'topics' ->> 0                         as topic_hash,
       CAST(log_entry ->> 'removed' AS BOOLEAN)            as removed
FROM expanded_logs
ORDER BY sequence, block_number desc;
