postgres:
	sudo docker run --name postgres16 --network banking-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine
createdb:
	sudo docker exec -it postgres16  createdb --username=root --owner=root simple_bank
dropdb:
	sudo docker exec -it postgres16 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
slqc:
	sqlc generate
server:
	go run main.go
network:
	docker network create bank-network
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Clayagiffeb/Simple_Bank/db/sqlc Store
.PHONY: network postgres createdb  dropdb migrateup migratedown  migrateup1 migratedown1 sqlc server mock