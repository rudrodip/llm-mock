APP_NAME=llmock
BIN_PATH=bin/$(APP_NAME)

build:
	@go build -o $(BIN_PATH) .

run: build
	@$(BIN_PATH)

clean:
	@rm -rf $(BIN_PATH)

test:
	@go test -v ./...