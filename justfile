start: 
    air

build:
	tailwindcss -i public/css/styles.css -o public/styles.css
	@templ generate view
	@go build -o bin/fullstackgo main.go 

test:
	@go test -v ./...
	
run: build
	@./bin/fullstackgo

tailwind:
	@tailwindcss -i views/css/styles.css -o public/styles.css --watch

templ:
	@templ generate -watch -proxy=http://localhost:8000

# migration: # add migration name at the end (ex: make migration create-cars-table)
# @migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@goose -dir migrations sqlite3 ./db.sqlite up
migrate-down:
	@goose -dir migrations sqlite3 ./db.sqlite down
db-status:
    @goose -dir migrations sqlite3 ./db.sqlite status

