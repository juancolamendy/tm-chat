CODE_ENTRY=cmd/main.go
APP_PORT=9090

.PHONY: run-server
run-server:
	go run $(CODE_ENTRY) -server

.PHONY: run-client
run-client:
	go run $(CODE_ENTRY)