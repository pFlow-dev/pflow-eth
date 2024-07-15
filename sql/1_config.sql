CREATE OR REPLACE FUNCTION config(p_name TEXT)
RETURNS TEXT AS $$
DECLARE
    config_data JSONB := $config$
        {
            "endpoint": "http://hardhat:8545",
            "service": "http://app:8080",
            "address": "0x5fbdb2315678afecb367f032d93f642f64180aa3"
        }
        $config$::jsonb;
BEGIN
    RETURN config_data->>p_name;
END; $$ LANGUAGE plpgsql;
