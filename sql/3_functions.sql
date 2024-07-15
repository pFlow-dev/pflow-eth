CREATE EXTENSION http;

CREATE OR REPLACE FUNCTION get_eth_transactions(addresses TEXT[], block_number INT) RETURNS TABLE(transaction_hash TEXT, transaction_details JSONB, logs JSONB) AS $$
DECLARE
    api_endpoint TEXT := (SELECT config('endpoint'));
    hex_block_number TEXT := '0x' || to_hex(block_number);
    block JSONB;
    transaction JSONB;
    transaction_hash TEXT;
    request_id INT;
    address TEXT;
BEGIN
    -- Get the next value from the sequence for the request ID
    request_id := nextval('request_id_seq');

    -- Get the block details
    SELECT content::jsonb INTO block
    FROM http_post(
        api_endpoint,
        '{
            "jsonrpc": "2.0",
            "method": "eth_getBlockByNumber",
            "params": ["' || hex_block_number || '", true],
            "id": ' || request_id || '
        }',
        'application/json'
    );

    -- Extract transactions involving any of the specific addresses
    FOR transaction IN SELECT * FROM jsonb_array_elements(block->'result'->'transactions')
    LOOP
        FOREACH address IN ARRAY addresses
        LOOP
            IF transaction->>'from' = address OR transaction->>'to' = address THEN
                transaction_hash := transaction->>'hash';

                -- Get the next value from the sequence for the request ID
                request_id := nextval('request_id_seq');

                -- Get transaction details
                SELECT content::jsonb INTO transaction
                FROM http_post(
                    api_endpoint,
                    '{
                        "jsonrpc": "2.0",
                        "method": "eth_getTransactionByHash",
                        "params": ["' || transaction_hash || '"],
                        "id": ' || request_id || '
                    }',
                    'application/json'
                );

                -- Get the next value from the sequence for the request ID
                request_id := nextval('request_id_seq');

                -- Get transaction receipt to fetch logs
                SELECT content::jsonb INTO logs
                FROM http_post(
                    api_endpoint,
                    '{
                        "jsonrpc": "2.0",
                        "method": "eth_getTransactionReceipt",
                        "params": ["' || transaction_hash || '"],
                        "id": ' || request_id || '
                    }',
                    'application/json'
                );

                -- Return the transaction details and logs
                RETURN QUERY SELECT transaction_hash, transaction, logs->'result'->'logs';
            END IF;
        END LOOP;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_eth_transactions(address TEXT, block_number INT) RETURNS TABLE(transaction_hash TEXT, transaction_details JSONB, logs JSONB) AS $$
BEGIN
    -- Call the modified get_eth_transactions function with the address wrapped in an array
    RETURN QUERY SELECT * FROM get_eth_transactions(ARRAY[address], block_number);
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_latest_block_number() RETURNS INT AS $$
DECLARE
    api_endpoint TEXT := (SELECT config('endpoint'));
    response JSONB;
    latest_block_number INT;
BEGIN
    SELECT content::jsonb INTO response
    FROM http_post(
            api_endpoint,
            '{
                "jsonrpc": "2.0",
                "method": "eth_blockNumber",
                "params": [],
                "id": 0
            }',
            'application/json'
         );

    latest_block_number = ('x' || lpad(ltrim(response->>'result', '0x'), 16, '0'))::bit(64)::bigint;

    RETURN latest_block_number;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sync_blocks() RETURNS VOID AS $$
DECLARE
    max_block_number INT;
    latest_block_number INT;
BEGIN
    SELECT MAX(block_number) INTO max_block_number FROM block_numbers;

    IF max_block_number IS NULL THEN
        max_block_number := 0;
    END IF;

    -- Get the latest block number using the get_latest_block_number function
    latest_block_number := get_latest_block_number();

    -- REVIEW: may want to do some number of max at a time
    WHILE max_block_number < latest_block_number LOOP
            max_block_number := max_block_number + 1;
            INSERT INTO block_numbers (block_number) VALUES (max_block_number);
        END LOOP;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_block_stats() RETURNS TABLE(highest_index INT, latest INT, behind INT) AS $$
DECLARE
    latest_block_number INT;
BEGIN
    -- Get the latest block number once
    latest_block_number := get_latest_block_number();

    RETURN QUERY
    SELECT
        max(block_number) AS highest_index,
        latest_block_number AS latest,
        latest_block_number - max(block_number) AS behind
    FROM
        block_numbers;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_block_by_number(block_number INT) RETURNS JSONB AS $$
DECLARE
    api_endpoint TEXT := (SELECT config('endpoint'));
    hex_block_number TEXT := '0x' || to_hex(block_number);
    response JSONB;
    request_id INT;
BEGIN
    request_id := nextval('request_id_seq');
    -- Get the block details from the Ethereum node API
    SELECT content::jsonb INTO response
    FROM http_post(
            api_endpoint,
            '{
                "jsonrpc": "2.0",
                "method": "eth_getBlockByNumber",
                "params": ["' || hex_block_number || '", true],
                "id": ' || request_id || '
            }',
            'application/json'
         );

    RETURN response;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION uint256_hex_to_numeric (val text) RETURNS numeric(78,0) AS $$
BEGIN
    -- Remove leading '0x' if present
    IF left(val, 2) = '0x' THEN
        val := substr(val, 3);
    END IF;

    RETURN trunc(val_a + val_b + val_c + val_d + val_e)
    FROM( SELECT
        concat('x00000000', substr(val, 1, 8))::bit(64)::int8::numeric(78,0) * 2^224::numeric(78,0) AS val_a,
        concat('x00', substr(val, 9, 14))::bit(64)::int8::numeric(78,0) * 2^168::numeric(78,0) AS val_b,
        concat('x00', substr(val, 23, 14))::bit(64)::int8::numeric(78,0) * 2^112::numeric(78,0) AS val_c,
        concat('x00', substr(val, 37, 14))::bit(64)::int8::numeric(78,0) * 2^56::numeric(78,0) AS val_d,
        concat('x00', substr(val, 51, 14))::bit(64)::int8::numeric(78,0) AS val_e
    FROM (SELECT val AS value) AS x) AS x;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_chain_id() RETURNS INT AS $$
DECLARE
    api_endpoint TEXT := (SELECT config('endpoint'));
    response JSONB;
    chain_id bigint;
BEGIN
    -- Get the chain ID from the Ethereum node API
    SELECT content::jsonb INTO response
    FROM http_post(
            api_endpoint,
            '{
                "jsonrpc": "2.0",
                "method": "eth_chainId",
                "params": [],
                "id": 0
            }',
            'application/json'
         );

    chain_id := ('x' || lpad(ltrim(response->>'result', '0x'), 16, '0'))::bit(64)::bigint;

    RETURN chain_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION signal(action int, scalar bigint) RETURNS JSONB AS $$
DECLARE
    request_id INT;
    api_endpoint TEXT;
    response JSONB;
    transaction_hash TEXT;
BEGIN
    -- Generate a new request ID
    request_id := nextval('request_id_seq');

    -- Append the request ID to the API endpoint URL
    api_endpoint := config('service') || '/v0/signal?action=' || action::text || '&scalar=' || scalar::text || '&nonce=' || request_id::text;

    -- Perform the HTTP GET request and store the response
    SELECT content::jsonb INTO response
    FROM http_get(api_endpoint);

    -- Extract transaction hash from the response (assuming the response structure includes a transaction hash)
    transaction_hash := response->'response'->>'transaction_hash';

    -- Insert the transaction hash and response into the sent_transactions table
    INSERT INTO sent_transactions (id, transaction_hash, response)
    VALUES (request_id, transaction_hash, response);

    -- Return the JSON response
    RETURN response;
END;
$$ LANGUAGE plpgsql;

-- select uint256_hex_to_numeric(signal->'response'->'eventLog'-> 0 ->> 'data')  created_sequence from signal(2, 1);


CREATE OR REPLACE FUNCTION signal_many(actions INT[], scalars BIGINT[]) RETURNS JSONB AS $$
DECLARE
    api_endpoint TEXT := (SELECT config('service'));
    request_id INT;
    response JSONB;
    actions_str TEXT;
    scalars_str TEXT;
    transaction_hash TEXT;
BEGIN
    -- Convert actions and scalars arrays to comma-separated strings
    SELECT string_agg(action::TEXT, ',') INTO actions_str FROM unnest(actions) AS action;
    SELECT string_agg(scalar::TEXT, ',') INTO scalars_str FROM unnest(scalars) AS scalar;

    -- Generate a new request ID
    request_id := nextval('request_id_seq');

    -- Make the HTTP POST request
    SELECT content::jsonb INTO response
    FROM http_post(
        api_endpoint || '/v0/signal?action=' || actions_str || '&scalar=' || scalars_str || '&nonce=' || request_id::TEXT,
        '',
        'application/json'
    );

    -- Extract transaction hash from the response (assuming the response structure includes a transaction hash)
    transaction_hash := response->'response'->>'transaction_hash';

    -- Insert the transaction hash and response into the sent_transactions table
    INSERT INTO sent_transactions ("id", "transaction_hash", "response")
    VALUES (request_id, transaction_hash, response);

    RETURN response;
END;
$$ LANGUAGE plpgsql;

-- select * from signal_many(ARRAY[2,2,2], ARRAY[1,2,3]);


CREATE OR REPLACE FUNCTION confirm_transactions() RETURNS JSONB AS $$
DECLARE
    matched_count INT := 0;
    unmatched_count INT := 0;
    total_sent_transactions INT := 0;
    report JSONB;
BEGIN
    -- Count total sent_transactions before deletion
    SELECT COUNT(*) INTO total_sent_transactions FROM sent_transactions;

    -- Delete sent_transactions that have a matching transaction_hash in the transactions table
    WITH deleted AS (
        DELETE FROM sent_transactions
        WHERE transaction_hash IN (SELECT t.transaction_hash FROM sent_transactions st join transactions t on t.transaction_hash = st.transaction_hash)
        RETURNING *
    )
    SELECT COUNT(*) INTO matched_count FROM deleted;

    -- Calculate unmatched transactions count
    unmatched_count := total_sent_transactions - matched_count;

    -- Prepare the report
    report := jsonb_build_object(
        'confirmed', matched_count,
        'pending', unmatched_count
    );

    RETURN report;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION truncate_node_data() RETURNS void AS $$
BEGIN
    EXECUTE 'TRUNCATE TABLE block_numbers CASCADE';

    EXECUTE 'TRUNCATE TABLE sent_transactions CASCADE';

    EXECUTE 'TRUNCATE TABLE transactions CASCADE';

    REFRESH MATERIALIZED VIEW transaction_logs_view;

    NOTIFY node_sync_channel, 'Node data truncated';
END;
$$ LANGUAGE plpgsql;
