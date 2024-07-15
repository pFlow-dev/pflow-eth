function get_contract_model() {
    contract_address=$1
    api_endpoint="http://localhost:8545/"
    call_declaration="0xb1a6afd3"
    call_model="0x0ad9d052"

    request='{
            "jsonrpc": "2.0",
            "method": "eth_call",
            "params": [{
                "to": "'"$contract_address"'",
                "data": "'"$call_model"'"
            }, "latest"],
            "id": 1
        }'

    echo $request;

    # Construct the JSON-RPC request
    response=$(curl -X POST \
        -H "Content-Type: application/json" \
        --data "${request}" \
        $api_endpoint)

    echo $response
}

get_contract_model 0x5fbdb2315678afecb367f032d93f642f64180aa3
