CREATE OR REPLACE FUNCTION node_sync() RETURNS jsonb AS $$
DECLARE
    max_block_number_before_sync INT;
    max_block_number_after_sync INT;
    highest_sequence INT;
    db_start_time TIMESTAMP;
    db_end_time TIMESTAMP;
    db_duration INTERVAL;
    unconfirmed_tx_count INT;
BEGIN
    -- Capture start time
    SELECT current_timestamp INTO db_start_time;

    -- Get the max block number before sync
    SELECT MAX(block_number) INTO max_block_number_before_sync FROM block_numbers;

    -- Perform the synchronization steps
    PERFORM sync_blocks();
    REFRESH MATERIALIZED VIEW transaction_logs_view;

    -- Get the max block number after sync
    SELECT MAX(block_number) INTO max_block_number_after_sync FROM block_numbers;

    -- Get the highest sequence from transaction_logs_view
    SELECT MAX(sequence) INTO highest_sequence FROM transaction_logs_view;

    -- Capture end time
    SELECT current_timestamp INTO db_end_time;

    SELECT COUNT(*) INTO unconfirmed_tx_count FROM sent_transactions;

    -- Calculate duration
    db_duration := db_end_time - db_start_time;

    -- Return JSON object with the max block numbers, highest sequence, and DB operation time
    RETURN jsonb_build_object(
        'unconfirmed_tx', unconfirmed_tx_count,
        'sequence', highest_sequence,
        'block', max_block_number_after_sync,
        'blocks_added', max_block_number_after_sync - max_block_number_before_sync,
        'time', db_start_time,
        'job_timer', db_end_time - db_start_time
    );
END;
$$ LANGUAGE plpgsql;

-- CREATE EXTENSION pg_cron;

-- schedule runs every minute
-- SELECT cron.schedule('*/1 * * * * *', 'SELECT refresh_and_insert()');

CREATE OR REPLACE FUNCTION wrap_node_sync() RETURNS TABLE(sync_data jsonb) AS $$
BEGIN
    -- Notify any listening clients that the node_sync has completed
    NOTIFY node_sync_channel, 'Node sync completed';
    RETURN QUERY SELECT node_sync() AS sync_data;
END;
$$ LANGUAGE plpgsql;

DROP MATERIALIZED VIEW node_sync_data_view;
CREATE MATERIALIZED VIEW node_sync_data_view AS
 SELECT * FROM wrap_node_sync();
