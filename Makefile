

run:
	go run ./cmd/app/main.go

tests:
	pytest ./tests/tests.py -v
	pytest ./tests/tests2.py -v

swag:
	swag init -g ./cmd/app/main.go --output ./docs