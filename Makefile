.PHONY: shadowify
shadowify:
	go run cmd/shadowify/main.go

migration-down:
	goose down

compose-up:
	docker-compose -f docker-compose.yml up -d

compose-down:
	docker-compose -f docker-compose.yml down