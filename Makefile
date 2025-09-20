.PHONY: build run


build:
go build -o bin/node ./cmd/node


run:
./bin/node --datadir ./data --rpc 8545


clean:
rm -rf bin data
