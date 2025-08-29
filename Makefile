AIR_VERSION ?= v1.62.0

.PHONY: live/templ
# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ:
	go tool templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

.PHONY: live/server
# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/air-verse/air@${AIR_VERSION} \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

.PHONY: live/tailwind
# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	npx --yes tailwindcss -i ./input.css -o ./assets/styles.css --minify --watch

.PHONY: live/esbuild
# run esbuild to generate the index.js bundle in watch mode.
live/esbuild:
	npx --yes esbuild js/index.ts --bundle --outdir=assets/ --watch=forever

.PHONY: live/sync_assets
# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	go run github.com/air-verse/air@${AIR_VERSION} \
	--build.cmd "go tool templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

.PHONY: live
# start all 5 watch processes in parallel.
live: 
	make -j5 live/templ live/server live/tailwind live/esbuild live/sync_assets
