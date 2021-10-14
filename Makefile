CODE_ENTRY=cmd/main.go
APP_PORT=9090

.PHONY: run
run:
	go run $(CODE_ENTRY)