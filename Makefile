# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif


MAIN_PATH = tmp/bin/main
SYNC_ASSETS_COMMAND =	@air \
--build.cmd "templ generate --notify-proxy" \
--build.bin "true" \
--build.delay "100" \
--build.exclude_dir "" \
--build.include_dir "public" \
--build.include_ext "js,css" \
--screen.clear_on_rebuild true \
--log.main_only true


# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
templ:
	@templ generate --watch --proxy="http://localhost$(HTTP_LISTEN_ADDR)" --open-browser=false

# run air to detect any go file changes to re-build and re-run the server.
server:
	@air \
	--build.cmd "go build --tags dev -o ${MAIN_PATH} ./cmd/" --build.bin "${MAIN_PATH}" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--screen.clear_on_rebuild true \
	--log.main_only true

# run tailwindcss to generate the styles.css bundle in watch mode.
watch-assets:
	@npx tailwindcss -i app/assets/app.css -o public/assets/styles.css --watch  

# run esbuild to generate the index.js bundle in watch mode.
watch-esbuild:
	@npx esbuild app/assets/index.js --bundle --outdir=public/assets --watch

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
sync_assets:
	${SYNC_ASSETS_COMMAND}

# start the application in development
dev:
	@make -j5 templ server watch-assets watch-esbuild sync_assets

# build the application for production. This will compile your app
# to a single binary with all its assets embedded.
build:
	@npx tailwindcss -i app/assets/app.css -o ./public/assets/styles.css --minify
	@npx esbuild app/assets/index.js --bundle --outdir=public/assets
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app_prod cmd/main.go
	@echo "compiled you application with all its assets to a single binary => bin/app_prod"

db-status:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(APP_DB_NAME) goose -dir=$(MIGRATION_DIR) status

db-reset:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(APP_DB_NAME) goose -dir=$(MIGRATION_DIR) reset

# db-down:
# 	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(APP_DB_NAME) goose -dir=$(MIGRATION_DIR) down

db-up:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(APP_DB_NAME) goose -dir=$(MIGRATION_DIR) up

db-mig-create:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(APP_DB_NAME) goose -dir=$(MIGRATION_DIR) create $(filter-out $@,$(MAKECMDGOALS)) sql

# db-seed:
# 	@go run app/scripts/seed/main.go
