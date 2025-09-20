# graphene-proto

Starter Graphene prototype (PoS + sharding skeleton) — minimal runnable code for local testing.

## Quickstart

1. Install Go 1.21+.
2. `git clone` this repo.
3. `make build` (produces `bin/node`).
4. `./bin/node --datadir ./data --rpc 8545` to start a single node.
5. Use JSON-RPC at http://localhost:8545/rpc with methods:
   - `Graphene.SendTx` (params: {from, to, amount})
   - `Graphene.GetBalance` (params: {address})
   - `Graphene.RegisterValidator` (params: {address, stake})
   - `Graphene.Delegate` (params: {delegator, validator, amount})

Example curl:

```bash
curl -s -X POST --data '{"method":"Graphene.SendTx","params":[{"from":"alice","to":"bob","amount":10}],"id":1}' http://localhost:8545/rpc


Running the chain:
go run cmd/gfn/main.go init
go run cmd/gfn/main.go start
This launches a simple 2-validator Graphene prototype producing blocks.
