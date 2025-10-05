DC := ./deployment/docker-compose.yaml
DC_TEST = ./deployment/docker-compose.test.yml
# 	 --progress=plain

.PHONY: test
test:
	docker compose -f $(DC_TEST) up --build --exit-code-from tests

.PHONY: run_serv
run_serv:
	docker compose -v -f $(DC) build 
	docker compose -v -f $(DC) run

.PHONY: run_app
run_app:
	docker compose -v -f $(DC) run --build code_processor_app

.PHONY: run_rabbit
run_rabbit:
	docker compose -v -f $(DC) run --build -d rabbitmq 

.PHONY: run_worker
run_worker:
	docker compose -v -f $(DC) run --build code_processor_worker 

.PHONY: down_app
down_app:
	docker compose -f $(DC) down -v code_processor_app

.PHONY: down_rabbit
down_rabbit:
	docker compose -f $(DC) down -v rabbitmq

.PHONY: down_all
down_all:
	docker compose -f $(DC) down -v
	docker compose -f $(DC_TEST) down -v

.PHONY: down_worker
down_worker:
	docker compose -f $(DC) down -v code_processor_worker	

run:
	go run ./cmd/app/main.go

.PHONY: tests
tests:
# 	pytest ./tests/tests.py -v
# 	pytest ./tests/tests2.py -v
	pytest ./tests/tests3.py -v

swag:
	swag init -g ./cmd/app/main.go --output ./docs