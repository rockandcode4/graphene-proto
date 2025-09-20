.PHONY: build run test docker clean

build:
	go build -o bin/node ./cmd/node

run:
	./bin/node --datadir ./data --rpc 8545

test:
	go test ./...

docker:
	docker-compose up --build

clean:
	rm -rf bin data
