.PHONY: createdb dropdb migrate-up migrate-down migrate-force sqlcddd test run mockgen migrate-up1 migrate-down1 initdb

createdb:
	docker exec -t simplebank-postgres createdb --username=guncv --owner=guncv simplebank

dropdb:
	docker exec -t simplebank-postgres dropdb simplebank

# Apply all up migrations
migrate-up:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simple_bank?sslmode=disable" --verbose up

# Roll back all migrations
migrate-down:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simple_bank?sslmode=disable" --verbose down

migrate-up1:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simple_bank?sslmode=disable" --verbose up 1

migrate-down1:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simple_bank?sslmode=disable" --verbose down 1

# Force reset migration to version 1
migrate-force:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simple_bank?sslmode=disable" force 1

migrate-version:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simple_bank?sslmode=disable" version

sqlc: 
	sqlc generate

mock: 
	mockgen -destination db/mock/store.go -package mockdb github.com/guncv/Simple-Bank/db/sqlc Store

test:
	go test -v -cover ./...

run:
	go run main.go

initdb:
	docker compose -f docker-compose.db.yaml up -d

	@echo "⏳ Waiting for database to be ready..."
	@until docker exec postgres12 pg_isready -U user -d simplebank; do sleep 1; done

	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simple_bank?sslmode=disable" --verbose up
