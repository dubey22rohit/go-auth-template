users:
	@go build -o bin/users ./users
	@./bin/users

chat:
	@go build -o bin/chat ./chat
	@./bin/chat

websocket:
	@go build -o bin/websocket ./pkg/websocket
	@./bin/websocket


.PHONY: users chat
