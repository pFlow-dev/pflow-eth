CREATE TABLE block_numbers (
    block_number INT PRIMARY KEY
);

CREATE TABLE transactions (
    transaction_hash TEXT PRIMARY KEY,
    transaction_details JSONB,
    logs JSONB,
    block_number INT
);
CREATE INDEX idx_transactions_transaction_hash ON transactions(transaction_hash);

CREATE TABLE transaction_states (
    sequence INT PRIMARY KEY,
    transaction_hash TEXT,
    state JSONB
);
CREATE INDEX idx_transaction_states_transaction_hash ON transaction_states(transaction_hash);

CREATE TABLE sent_transactions (
    id INT PRIMARY KEY,
    transaction_hash TEXT,
    response JSONB
);
CREATE INDEX idx_sent_transactions_transaction_hash ON sent_transactions(transaction_hash);