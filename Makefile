format:
	gofumpt -l -w .
run:
	gofumpt -l -w .
	go run .
build:
	gofumpt -l -w .
	go build -o own-redis .
remove:
	rm -rf own-redis