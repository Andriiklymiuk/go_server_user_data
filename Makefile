.PHONY: run build dev triggerAll up down generateDbQueries

run:
	go run .
build:
	go build -o build
dev: 
	./run.sh
up:
	go run . --migrateUp
down:
	go run . --migrateDown
generateDbQueries:
	sqlc generate
	
triggerAll:
	tc --col all