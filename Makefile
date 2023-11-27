.DEFAULT_GOAL = help

MAIN_API = "cmd/api/main.go"
MAIN_API_BIN = "bin/api"

CURRENT_TIME = $(shell date --iso-8601=seconds)

build: ## Builds the cmd/api application.
	@echo 'Building cmd/api...'
	go build -ldflags='-s -X main.buildTime=${CURRENT_TIME}' -o=$(MAIN_API_BIN) $(MAIN_API)

unit-tests: ## Run unit tests in verbose mode.
	@echo 'Starting tests...'
	go test -v -cover -race ./...

swag-gen: ## Generate swagger files.
	@echo 'Generating swagger files...'
	swag init -q -g $(MAIN_API) -o docs/swagger

kubectl-config-context: ## Set the current-context in a kubeconfig file.
	kubectl config use-context arn:aws:eks:eu-north-1:425727356824:cluster/packer

kubectl-apply-conf: ## Apply the new configuration to the RBAC configuration of the Amazon EKS cluster.
	kubectl apply -f eks/aws-auth.yml

kubectl-apply-deploy: ## Apply the new deployment.
	kubectl apply -f eks/deployment.yaml

kubectl-apply-srvc: ## Apply the new service.
	kubectl apply -f eks/service.yaml

kubectl-service: ## List all services in the namespace.
	kubectl get service

kubectl-pods: ## List all pods in the namespace.
	kubectl get pods

help: ## Prints this message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build unit-tests swag-gen kubectl-config-contextkubectl-apply-conf kubectl-apply-deploy  kubectl-service \
 	kubectl-pods kubectl-apply-srvc help
