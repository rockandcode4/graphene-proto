#!/usr/bin/env bash
set -e


# quick localnet script (3 nodes) for testing
BIN=./bin/node
if [ ! -f "$BIN" ]; then
echo "build first: make build"
exit 1
fi


mkdir -p localnet/node1 localnet/node2 localnet/node3


# st
