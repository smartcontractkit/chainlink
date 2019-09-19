REPO := smartcontract/chainlink
DOCKER_TAG ?= latest
TAGGED_REPO := $(REPO):$(DOCKER_TAG)

.PHONY: dockerpush
dockerpush: ## Push the docker image to dockerhub
	docker push $(TAGGED_REPO)
