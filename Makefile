.PHONY: shadowify
run:
	go run cmd/shadowify/main.go

compose-up:
	docker-compose -f docker-compose.yml up -d

compose-down:
	docker-compose -f docker-compose.yml down