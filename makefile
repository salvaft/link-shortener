build: 
	tailwindcss -i views/css/styles.css -o public/styles.css
	templ generate views
	go build -o bin/main . 
	cd js && pnpm i && pnpm run build
