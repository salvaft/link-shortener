## Requirements
- tailwindcss
- templ
- just
- go 1.22
- pnpm (for conveniency)
- goose (SQL migration query runner)

In MacOS or Linux run:
```
brew install goose tailwindcss just
go install github.com/a-h/templ/cmd/templ@latest
```

## Startup
```sh
goose sqlite3 ./db.sqlite3 --dir migrations
go mod tidy
just run
```
