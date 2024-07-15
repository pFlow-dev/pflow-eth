#!/bin/bash

API="http://localhost:8545"

BLOCK=$(curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' -H "Content-Type: application/json" $API)

if [ -z "$BLOCK" ]; then
  echo "Error: Could not get the latest block number"
  exit 1
fi

echo $BLOCK
