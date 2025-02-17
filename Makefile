.PHONY: createdb dropdb migrate-up migrate-down migrate-force sqlcddd test run mockgen migrate-up1 migrate-down1

createdb:
	docker exec -t simplebank-postgres createdb --username=guncv --owner=guncv simplebank

dropdb:
	docker exec -t simplebank-postgres dropdb simplebank

# Apply all up migrations
migrate-up:
	migrate -path db/migration -database "postgresql://guncv:17554Guncv_26042003@localhost:5432/simplebank?sslmode=disable" --verbose up

# Roll back all migrations
migrate-down:
	migrate -path db/migration -database "postgresql://guncv:17554Guncv_26042003@localhost:5432/simplebank?sslmode=disable" --verbose down

migrate-up1:
	migrate -path db/migration -database "postgresql://guncv:17554Guncv_26042003@localhost:5432/simplebank?sslmode=disable" --verbose up 1

migrate-down1:
	migrate -path db/migration -database "postgresql://guncv:17554Guncv_26042003@localhost:5432/simplebank?sslmode=disable" --verbose down 1

# Force reset migration to version 1
migrate-force:
	migrate -path db/migration -database "postgresql://guncv:17554Guncv_26042003@localhost:5432/simplebank?sslmode=disable" force 1

migrate-version:
	migrate -path db/migration -database "postgresql://guncv:17554Guncv_26042003@localhost:5432/simplebank?sslmode=disable" version

sqlc: 
	sqlc generate

mock: 
	mockgen -destination db/mock/store.go -package mockdb github.com/guncv/Simple-Bank/db/sqlc Store

test:
	go test -v -cover ./...

run:
	go run main.go