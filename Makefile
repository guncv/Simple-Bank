.PHONY: createdb dropdb migrate-up migrate-down migrate-force sqlcddd test

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

# Force reset migration to version 1
migrate-force:
	migrate -path db/migration -database "postgresql://guncv:17554Guncv_26042003@localhost:5432/simplebank?sslmode=disable" force 1

migrate-version:
	migrate -path db/migration -database "postgresql://guncv:17554Guncv_26042003@localhost:5432/simplebank?sslmode=disable" version

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...