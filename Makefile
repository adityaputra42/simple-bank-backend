postgres:
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:16-alpine

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
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple-bank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go simple-bank/worker TaskDistributor


proto:	
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
  proto/*.proto
	statik -src=./doc/swagger -dest=./doc

redis:
	docker run --name redis -p 6379:6379 -d redis:8-alpine

.PHONY:postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 proto redis