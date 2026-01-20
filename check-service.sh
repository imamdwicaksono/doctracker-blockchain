# Redis
redis-cli ping

# Ethereum RPC
curl -X POST http://127.0.0.1:8545 \
-H "Content-Type: application/json" \
-d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Go app
curl http://localhost:8080
