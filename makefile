## run the app

start:
	docker compose down
	@echo "Deploying app on docker ...."
	docker compose  -f docker-compose.yml  -f docker-compose.production.yml up --build

stop:
	@echo "Shutting down app.."
	docker compose -f docker-compose.yml -f docker-compose.production.yml rm server
	docker compose down

