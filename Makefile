start:
	@go run cmd/server/main.go

mod:
	go mod tidy
	go mod vendor

specs: ## Generate swagger specs
	HOST=$(HOST) sh scripts/specs-gen.sh

swarm:
	docker stack deploy --compose-file ./.docker/swarm/docker-compose.dev.yml go-clean-arch

docker-config:
	docker config create go-clean-arch.conf config.dev.yml

build-image:
	docker build -t lhhoangit/go-clean-arch .

lint: 
	golangci-lint version
	golangci-lint run -v -c golangci.yml ./...