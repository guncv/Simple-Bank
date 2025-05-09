.PHONY: createdb dropdb migrate-up migrate-down migrate-force sqlcddd test reset build remove-img mockgen migrate-up1 migrate-down1 initdb network proto

createdb:
	docker exec -t simplebank-postgres createdb --username=guncv --owner=guncv simplebank

dropdb:
	docker exec -t simplebank-postgres dropdb simplebank

# Apply all up migrations
migrate-up:
	migrate -path db/migration -database "postgresql://user:password@postgres12:5432/simple_bank?sslmode=disable" --verbose up

migrate-ci-up:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simple_bank?sslmode=disable" --verbose up

# Roll back all migrations
migrate-down:
	migrate -path db/migration -database "postgresql://user:password@postgres12:5432/simple_bank?sslmode=disable" --verbose down

migrate-up1:
	migrate -path db/migration -database "postgresql://user:password@postgres12:5432/simple_bank?sslmode=disable" --verbose up 1

migrate-down1:
	migrate -path db/migration -database "postgresql://user:password@postgres12:5432/simple_bank?sslmode=disable" --verbose down 1

# Force reset migration to version 1
migrate-force:
	migrate -path db/migration -database "postgresql://user:password@postgres12:5432/simple_bank?sslmode=disable" force 1

migrate-version:
	migrate -path db/migration -database "postgresql://user:password@postgres12:5432/simple_bank?sslmode=disable" version

network:
	docker network create simple-bank-network

sqlc: 
	sqlc generate

mock: 
	mockgen -destination db/mock/store.go -package mockdb github.com/guncv/Simple-Bank/db/sqlc Store

test:
	go test -v -cover ./...

build:
	docker compose build

reset:
	docker compose down && docker compose build --no-cache && docker compose up

proto: 
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simplebank \
	proto/*.proto
	statik -src=./docs/swagger -dest=./docs

evans: 
	evans \
		--host localhost \
		--port 9090 \
		--path proto \
		--proto service_simple_bank.proto \
		--proto rpc_create_user.proto \
		--proto rpc_login_user.proto \
		--proto user.proto \
		repl




