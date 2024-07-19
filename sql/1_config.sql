CREATE OR REPLACE FUNCTION config(p_name TEXT)
    RETURNS TEXT AS
$$
DECLARE
    p_network TEXT := current_schema();
    config_data JSONB := $config$
    {
      "hardhat": {
        "service": "http://app:8080",
        "endpoint": "http://hardhat:8545",
        "address": "0x5fbdb2315678afecb367f032d93f642f64180aa3",
        "testnet": true
      },
      "sepolia_optimism": {
        "endpoint": "https://optimism-sepolia.blastapi.io/4fd309c5-6555-4052-99d5-3a6b646f14b4",
        "address": "0x5fbdb2315678afecb367f032d93f642f64180aa3",
        "testnet": true
      },
      "optimism": {
        "endpoint": "https://optimism-mainnet.blastapi.io/4fd309c5-6555-4052-99d5-3a6b646f14b4",
        "address": "0x5fbdb2315678afecb367f032d93f642f64180aa3",
        "testnet": false
      }
    }
    $config$::jsonb;
BEGIN
    if p_name = 'schema' then
        RETURN p_network;
    end if;
    RETURN config_data -> p_network ->> p_name;
END;
$$ LANGUAGE plpgsql;
