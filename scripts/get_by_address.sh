#!/bin/bash

ADDRESS="0x7f1ed3d3aac8903f869eeb32182265dc34106353"
BLOCK_NUMBER=9934458
HEX_BLOCK_NUMBER=$(printf "0x%x" $BLOCK_NUMBER)
API="https://dark-patient-fog.optimism-sepolia.quiknode.pro/25f53f89a1e74fad042f899a2a46643f728ab98f/"

# Get the block details
BLOCK=$(curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["'"$HEX_BLOCK_NUMBER"'", true],"id":1}' -H "Content-Type: application/json" $API)

# Extract transactions involving the specific address
TRANSACTION_HASHES=$(echo $BLOCK | jq -r '.result.transactions[] | select(.from == "'$ADDRESS'" or .to == "'$ADDRESS'") | .hash')

# Function to get transaction receipt by hash
get_transaction_receipt() {
    local tx_hash=$1
    curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["'"$tx_hash"'"],"id":1}' -H "Content-Type: application/json" $API
}

# Loop through each transaction hash and get details and logs
echo "Detailed transactions involving $ADDRESS in block $BLOCK_NUMBER:"
for tx_hash in $TRANSACTION_HASHES; do
    echo "Transaction Hash: $tx_hash"
    TRANSACTION_DETAILS=$(curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["'"$tx_hash"'"],"id":1}' -H "Content-Type: application/json" $API)
    echo $TRANSACTION_DETAILS | jq
    
    echo "Events (Logs):"
    TRANSACTION_RECEIPT=$(get_transaction_receipt $tx_hash)
    echo $TRANSACTION_RECEIPT | jq '.result.logs'
    echo "----------------------------------------"
done
