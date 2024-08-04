.PHONY: docker-run
docker-run: docker-build
	docker run kleos

.PHONY: docker-build
docker-build:
	docker build -t kleos -f Dockerfile.kleos .

.PHONY: migration-up
migration-up:
	migrate -path db_migration/migrations/ -database "${POSTGRES_DSN}" -verbose up

.PHONY: migration-down
migration-down:
	migrate -path db_migration/migrations/ -database "${POSTGRES_DSN}" -verbose down

.PHONY: build
build:
	go mod tidy && CGO_ENABLED=0 go build -ldflags="-s -w" -o kleos main.go