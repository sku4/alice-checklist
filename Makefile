.PHONY:

export GOOS=linux
build:
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

swagger:
	swag init -g ./cmd/main.go

init-project-k8s:
	@./scripts/init-project-k8s.sh alice prod

helm-install:
	helm upgrade --install "alice-checklist" .helm --namespace=alice-prod

helm-install-local:
	helm upgrade --install "alice-checklist" .helm \
		--namespace=alice-prod \
		-f ./.helm/values-local.yaml

helm-template:
	helm template --name-template="alice-checklist" --namespace=alice-prod -f .helm/values-local.yaml .helm > .helm/helm.txt
