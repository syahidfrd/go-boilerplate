.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: build
build:
	go build -o /tmp/bin/app main.go

.PHONY: run
run: build
	/tmp/bin/app $(bin)

.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/app $(bin)" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

.PHONY: test/unit
test/unit:
	go test ./... -race -coverprofile=./coverage.out

.PHONY: test/integration
test/integration:
	go test -tags=integration ./... -v -timeout 5m

.PHONY: test/all
test/all: test test/integration
