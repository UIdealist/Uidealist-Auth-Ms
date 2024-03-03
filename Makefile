# Common variables for this application
APP_NAME = uidealist-auth-ms
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/platform/migrations
DATABASE_URL=

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

# Migrations commands (Define DATABASE_URL)
migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

# Main Application docker commands
NETWORK_NAME=dev-network-${APP_NAME}
docker.network:
	docker network create -d bridge ${NETWORK_NAME}

docker.network.destroy:
	docker network rm ${NETWORK_NAME}

docker.fiber.build:
	docker build -t ${APP_NAME} .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name ${APP_NAME} \
		--network ${NETWORK_NAME} \
		--env-file .env \
		-p 5000:5000 \
		${APP_NAME}

# Database for test purposes
docker.postgres:
	docker run --rm -d \
		--name ${APP_NAME}-postgres \
		--network ${NETWORK_NAME} \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=postgres \
		-v ${HOME}/dev-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

# Redis cache for test purposes
docker.redis:
	docker run --rm -d \
		--name ${APP_NAME}-redis \
		--network ${NETWORK_NAME} \
		-p 6379:6379 \
		redis

# Stop commands for containers
docker.stop.fiber:
	docker stop ${APP_NAME}

docker.stop.postgres:
	docker stop ${APP_NAME}-postgres

docker.stop.redis:
	docker stop ${APP_NAME}-redis

docker.run: docker.postgres docker.redis docker.fiber
docker.stop: docker.stop.fiber docker.stop.postgres docker.stop.redis