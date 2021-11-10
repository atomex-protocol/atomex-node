-include .env
export $(shell sed 's/=.*//' .env)

tower:
	cd cmd/watch_tower && go run .

tower-test:
	cd cmd/watch_tower && go run . -c config.test.yml

market-maker:
	cd cmd/market_maker && go run . -c ../../configs/test

test:
	go test ./...

lint:
	golangci-lint run

up:
	docker-compose up -d --build watch_tower