DB_URL=postgresql://postgres:cong1234@localhost:5432/simple-bank-app?sslmode=disable

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple-bank-app

dropdb:
	docker exec -it postgres dropdb simple-bank-app	

postgres:
	docker run --name postgres --network simple-bank-app -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate	

test:
	go test -v -cover -short ./...

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/Sotatek-CongNguyen/simple-bank-practice/db/sqlc Store

proto: 
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans: 
	evans --host localhost --port 9090 -r repl

.PHONY: postgres createdb dropdb postgres migrateup migratedown sqlc test server mockgen proto evans
