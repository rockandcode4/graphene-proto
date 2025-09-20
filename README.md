Graphene-proto — prototype of Graphene chain (PoS + sharding skeleton)

Quickstart:
1. go build ./cmd/node
2. ./node --datadir ./node1 --rpc 8545

Use HTTP JSON-RPC at /rpc with method Graphene.SendTx and Graphene.GetBalance.
Example (curl):
curl -s -X POST --data '{"method":"Graphene.SendTx","params":[{"from":"alice","to":"bob","amount":100}],"id":1}' http://localhost:8545/rpc
