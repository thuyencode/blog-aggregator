ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export GOOSE_DBSTRING="postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
endif

test:
	go test -json ./... | gotestfmt
fmt:
	golangci-lint fmt
lint:
	golangci-lint run
fix:
	golangci-lint run --fix
clean:
	go clean

db-migrate-up:
	goose up
db-migrate-down:
	goose down
db-migrate-status:
	goose status
