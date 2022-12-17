.PHONY: run-server
run-server:
	go run server/server.go

.PHONY: run-client
run-client:
	go run client/client.go
