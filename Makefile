install:
	brew install grpcurl
	go install github.com/spf13/cobra-cli@latest

test:
	go test -race ./...

bench:
	go test ./... -bench=.

serve:
	go run main.go serve