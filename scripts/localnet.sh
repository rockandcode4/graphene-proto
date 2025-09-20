#!/usr/bin/env bash
set -e

make build

mkdir -p localnet/node1 localnet/node2 localnet/node3

# start nodes in background
./bin/node --datadir ./localnet/node1 --rpc 8545 &
./bin/node --datadir ./localnet/node2 --rpc 8546 &
./bin/node --datadir ./localnet/node3 --rpc 8547 &

echo "Nodes started. Use /rpc endpoint on ports 8545/8546/8547"
