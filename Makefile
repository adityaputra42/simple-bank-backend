postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -path db/migration -verbose up

migratedown:
	migrate -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -path db/migration -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY:postgres createdb dropdb migrateup migratedown sqlc test