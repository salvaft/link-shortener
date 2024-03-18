## Requirements
- tailwindcss
- templ
- make
- go 1.22
- pnpm (for conveniency)
- goose (SQL migration query runner)


## Startup
```sh
goose sqlite3 ./db.sqlite3 --dir migrations
go mod tidy
just run
```
