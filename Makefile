.PHONY:

export GOOS=linux
build:
	swag init -g ./cmd/main.go
	go build -o ./.bin/app ./cmd/main.go

run: build
	docker-compose up --remove-orphans --build server

test:
	go test ./... -coverprofile cover.out

test-coverage:
	go tool cover -func cover.out | grep total | awk '{print $3}'

build-image:
	docker build -t sku4/alice-checklist:v1.0.0 .

start-container:
	docker run \
		-v $(CURDIR)/db:/root/db \
		-v $(CURDIR)/configs/googlekeep:/root/configs/googlekeep \
		--env-file .env \
		-p 8000:8000 \
		sku4/alice-checklist:v1.0.0
