current_dir = $(shell pwd)

DOCKER_REGISTRY_SERVER=registry.gitlab.com/balconygames/analytics
DOCKER_REGISTRY_DASHBOARD=registry.gitlab.com/balconygames/analytics/dashboard
DOCKER_COMPOSE=pipenv run docker-compose
ANSIBLE=$(DOCKER_COMPOSE) -f deploy/ansible/docker-compose.deploy.yml run

init:
	pip install --user pipenv
	pipenv install
	cd deploy/ansible/ ; pipenv install ; cd ../../
.PHONY: init

up:
	${DOCKER_COMPOSE} up -d
.PHONY: up

stop:
	${DOCKER_COMPOSE} stop
.PHONY: stop

rm:
	${DOCKER_COMPOSE} rm -y
.PHONY: rm

db:
	$(MAKE) -C modules/primary db
	$(MAKE) -C modules/auth db
	$(MAKE) -C modules/leaderboard db
.PHONY: db

lint:
	docker run --rm -v $(current_dir):/app -w /app golangci/golangci-lint:v1.24.0 golangci-lint run -v ./...
.PHONY: lint

release: release.server release.dashboard
.PHONY: release

release.server:
	docker build -f Dockerfile.release -t $(DOCKER_REGISTRY_SERVER):latest .
	docker push $(DOCKER_REGISTRY_SERVER):latest
.PHONY: release.server

release.dashboard:
	cd web/ ; docker build -f Dockerfile.release -t $(DOCKER_REGISTRY_DASHBOARD):latest .
	cd web/ ; docker push $(DOCKER_REGISTRY_DASHBOARD):latest
.PHONY: release.dashboard

console:
	ssh root@$(PROD_SERVER)
.PHONY: console

deploy: release deploy.release
.PHONY: deploy

deploy.provision:
	$(ANSIBLE) provision
.PHONY: deploy.provision

deploy.release:
	$(ANSIBLE) release
.PHONY: deploy.release

