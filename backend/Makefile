# Install CompileDaemon for local dev: go install github.com/githubnemo/CompileDaemon@latest

.PHONY: dev tidy run-dev run-prod

dev:
	CompileDaemon \
		-build="go build -o tmp/main ." \
		-command="./tmp/main -env=development" \
		-directory="." \
		-include="*.go,*.env.*"

tidy:
	go mod tidy

run-dev:
	go run . -env=development

run-prod:
	go run . -env=production
