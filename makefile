## run the app

prod-run:
	sudo docker compose -p lenspix down
	@echo "Deploying on docker in ec2 instance"
	sudo docker compose  -f docker-compose.yml  -f docker-compose.production.yml up --build


run:
	docker compose -p lenspix down
	@echo "Deploying app on docker ...."
	docker compose  -f docker-compose.yml  -f docker-compose.production.yml up --build

stop:
	@echo "Shutting down app.."
	docker compose -f docker-compose.yml -f docker-compose.production.yml rm server
	docker compose down

