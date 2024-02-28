postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -path db/migration -verbose up

migrateup1:
	migrate -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -path db/migration -verbose up 1

migratedown:
	migrate -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -path db/migration -verbose down

migratedown1:
	migrate -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -path db/migration -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple-bank/db/sqlc Store

.PHONY:postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1