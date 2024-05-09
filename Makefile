build:
	@go build -o bin/api

run: build
	@./bin/api

test:
	@go test -v ./...

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 8080:8080 api

seed:
	@go run scripts/seed.go
