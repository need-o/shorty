.PHONY: app
build:
	go build 
run:
	go build; ./shorty
migrate_up:
	migrate -database sqlite3://shorty.db -path migrations up 
migrate_down:
	migrate -database sqlite3://shorty.db -path migrations down

.DEFAULT_GOAL := build