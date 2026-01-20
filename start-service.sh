#!/bin/bash

LOG_DIR="./logs"
mkdir -p $LOG_DIR

echo "🚀 Starting Redis..."
redis-server > $LOG_DIR/redis.log 2>&1 &
REDIS_PID=$!

echo "⛓️  Starting Anvil..."
anvil --host 127.0.0.1 --port 8545 > $LOG_DIR/anvil.log 2>&1 &
ANVIL_PID=$!

echo "🐹 Starting Go app..."
go run main.go > $LOG_DIR/go.log 2>&1 &
GO_PID=$!

echo "✅ Running in background"
echo "Redis PID : $REDIS_PID"
echo "Anvil PID : $ANVIL_PID"
echo "Go PID    : $GO_PID"

echo "$REDIS_PID $ANVIL_PID $GO_PID" > run.pid
