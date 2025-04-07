build:
		@go build -o bin/dfs_go

run: build
		@./bin/dfs_go

test:
		@go test -v ./...