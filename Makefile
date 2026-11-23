.PHONY: gen db-reset migrate-up run tools deploy build-bin migrate-remote seeder-master seed-remote-master

build:
	docker compose build

build-nc:
	docker compose build --no-cache

up:
# 	@make set-up-githooks
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart

ps:
	docker compose ps

# --- Staging Server Commands ---
build-staging:
	docker compose -f docker-compose.staging.yml build

build-staging-nc:
	docker compose -f docker-compose.staging.yml build --no-cache

up-staging:
	docker compose -f docker-compose.staging.yml up -d

down-staging:
	docker compose -f docker-compose.staging.yml down

restart-staging:
	docker compose -f docker-compose.staging.yml restart

ps-staging:
	docker compose -f docker-compose.staging.yml ps

logs-staging:
	docker compose -f docker-compose.staging.yml logs -f platform-api

migrate-staging-up:
	docker compose -f docker-compose.staging.yml exec platform-api migrate -path db/migrations -database "mysql://root:root@tcp(platform-db)/platform_db" up

set-up-githooks:
	git config --local core.hooksPath .githooks
	chmod -R +x .githooks/


# Generate backend + frontend (Alias for codegen)
gen:
	@make codegen

# Install requirement tools
tools:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/air-verse/air@v1.61.1
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

db-reset:
	docker compose exec -T platform-db mysql -uroot -proot -e "DROP DATABASE IF EXISTS platform_db; CREATE DATABASE platform_db;"


# run:
# 	air -c .air.toml

air:
	docker compose exec platform-api air -c .air.toml

lint:
	docker compose exec platform-api golangci-lint run

wire:
	docker compose exec platform-api wire gen ./internal/di/wire.go
	docker compose exec platform-api wire gen -output_file_prefix=test_ ./internal/di/wire-test.go

# OpenAPI Bundle
bundle:
	mkdir -p src/generated/api src/generated/openapi/openapi
	docker compose exec platform-node openapi-generator-cli generate -c platform-openapitools.json

codegen:
	@make bundle
	docker compose exec platform-api oapi-codegen -config config/server.yaml generated/openapi/openapi/openapi.yaml
	docker compose exec platform-api oapi-codegen -config config/models.yaml generated/openapi/openapi/openapi.yaml
# 	@make codegen-ts

codegen-ts:
	docker compose exec platform-node npx openapi2aspida -i=./generated/openapi/openapi/openapi.yaml -o=./generated/backend-api
	@if [ -d "../project-frontend" ]; then \
		echo "Moving backend-api files to ../project-frontend..."; \
		cp -r src/generated/backend-api/@types/ ../project-frontend/src/generated/backend-api/@types/ 2>/dev/null || true; \
	fi
	docker compose exec platform-node rm -rf generated/backend-api

#--------------------- migration start
migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migration name=<migration_name>"; \
		exit 1; \
	fi
	docker compose exec platform-api migrate create -ext sql -dir db/migrations $(name)

migrate-up:
	docker compose exec platform-api migrate -path db/migrations -database "mysql://root:root@tcp(platform-db)/platform_db" up

migrate-down:
	docker compose exec platform-api migrate -path db/migrations -database "mysql://root:root@tcp(platform-db)/platform_db" down -all

migrate-refresh:
	@make migrate-down
	@make migrate-up

migrate-fresh:
	docker compose exec platform-db mysql --user=root --password=root -e 'DROP DATABASE IF EXISTS `platform_db`;'
	docker compose exec platform-db mysql --user=root --password=root -e 'CREATE DATABASE `platform_db`;'
	@make migrate-up

seeder:
	docker compose exec platform-api go run db/seed/seed.go

seeder-master:
	@if [ -z "$(table)" ]; then \
		docker compose exec platform-api go run db/seed/seed.go --master; \
	else \
		docker compose exec platform-api go run db/seed/seed.go --master=$(table); \
	fi

#--------------------- test & test-db start
test:
	docker compose exec platform-api go test -v ./...

test-controller:
	@if [ -z "$(controller)" ]; then \
		docker compose exec platform-api go test -v ./internal/controller_test/...; \
	else \
		docker compose exec platform-api go test -v ./internal/controller_test/$(controller)/...; \
	fi

test-migrate-up:
	docker compose exec platform-api migrate -path db/migrations -database "mysql://root:root@tcp(platform-db-test)/platform_db_test" up

test-migrate-down:
	docker compose exec platform-api migrate -path db/migrations -database "mysql://root:root@tcp(platform-db-test)/platform_db_test" down -all

test-migrate-refresh:
	@make test-migrate-down
	@make test-migrate-up

test-migrate-fresh:
	docker compose exec platform-db-test mysql --user=root --password=root -e 'DROP DATABASE IF EXISTS `platform_db_test`;'
	docker compose exec platform-db-test mysql --user=root --password=root -e 'CREATE DATABASE `platform_db_test`;'
	@make test-migrate-up

test-reset-db:
	docker compose exec platform-db-test mysql --user=root --password=root -e 'DROP DATABASE IF EXISTS `platform_db_test`; CREATE DATABASE `platform_db_test`;'
	@make test-migrate-up

db-test:
	docker compose exec platform-db-test mysql -uroot -proot platform_db_test

# ==============================================================================
# Remote Server Build & Deploy (Staging / Production)
# ==============================================================================
SSH_USER      ?= ymsdwuqshosting
SSH_HOST      ?= 202.92.5.66
SSH_PORT      ?= 65333
SSH_KEY       ?=
REMOTE_GO_DIR ?= ~/platform-api
GO_BIN_NAME   ?= server
GO_BIN        := bin/$(GO_BIN_NAME)

SSH_KEY_OPT   := $(if $(SSH_KEY),-i $(SSH_KEY),)
SSH_CMD       := ssh $(SSH_KEY_OPT) -p $(SSH_PORT)

test-ssh:
	@echo "Kiểm tra kết nối SSH tới $(SSH_USER)@$(SSH_HOST):$(SSH_PORT)..."
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "echo 'SSH connected successfully as user:' \$$(whoami)"

build-bin:
	@echo "------------------------------------------------------------------------"
	@echo "===> Compile Go binary cho Linux amd64 (platform-go)..."
	@echo "------------------------------------------------------------------------"
	@mkdir -p bin
	@cd src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../bin/$(GO_BIN_NAME) ./cmd/server
	@chmod +x $(GO_BIN)
	@echo "===> Hoàn tất build Go binary: $(GO_BIN)"

deploy: build-bin
	@echo "------------------------------------------------------------------------"
	@echo "===> Deploying Go Backend lên $(SSH_USER)@$(SSH_HOST):$(REMOTE_GO_DIR)..."
	@echo "------------------------------------------------------------------------"
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "mkdir -p $(REMOTE_GO_DIR)/db"
	@echo "Dừng server cũ trước khi cập nhật binary..."
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "test -f $(REMOTE_GO_DIR)/service.sh && bash $(REMOTE_GO_DIR)/service.sh stop" || true
	@echo "Copying Go binary..."
	@scp $(SSH_KEY_OPT) -P $(SSH_PORT) $(GO_BIN) $(SSH_USER)@$(SSH_HOST):$(REMOTE_GO_DIR)/$(GO_BIN_NAME).new
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "mv -f $(REMOTE_GO_DIR)/$(GO_BIN_NAME).new $(REMOTE_GO_DIR)/$(GO_BIN_NAME) && chmod +x $(REMOTE_GO_DIR)/$(GO_BIN_NAME)"
	@if [ -d "src/db/migrations" ]; then \
		echo "Copying database migrations..."; \
		scp $(SSH_KEY_OPT) -P $(SSH_PORT) -r src/db/migrations $(SSH_USER)@$(SSH_HOST):$(REMOTE_GO_DIR)/db/; \
	fi
	@echo "Syncing service.sh..."
	@scp $(SSH_KEY_OPT) -P $(SSH_PORT) service.sh $(SSH_USER)@$(SSH_HOST):$(REMOTE_GO_DIR)/service.sh
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "chmod +x $(REMOTE_GO_DIR)/service.sh"
	@echo "Kiểm tra file .env trên server..."
	@if [ -f "src/.env.staging" ]; then \
		$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "test -f $(REMOTE_GO_DIR)/.env" || \
		(echo "Khởi tạo $(REMOTE_GO_DIR)/.env từ .env.staging..." && scp $(SSH_KEY_OPT) -P $(SSH_PORT) src/.env.staging $(SSH_USER)@$(SSH_HOST):$(REMOTE_GO_DIR)/.env); \
	fi
	@echo "Chạy Database Migrations trên server..."
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "bash $(REMOTE_GO_DIR)/service.sh migrate" || true
	@echo "Chạy Seeder Master Data trên server..."
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "bash $(REMOTE_GO_DIR)/service.sh seed-master" || true
	@echo "Starting backend trên server..."
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "bash $(REMOTE_GO_DIR)/service.sh restart" || true
	@echo "===> Deploy Go Backend thành công!"

migrate-remote:
	@echo "------------------------------------------------------------------------"
	@echo "===> [MIGRATE] Chạy Database Migrations trên $(SSH_USER)@$(SSH_HOST)..."
	@echo "------------------------------------------------------------------------"
	@if [ -d "src/db/migrations" ]; then \
		echo "Syncing database migrations..."; \
		scp $(SSH_KEY_OPT) -P $(SSH_PORT) -r src/db/migrations $(SSH_USER)@$(SSH_HOST):$(REMOTE_GO_DIR)/db/; \
	fi
	@scp $(SSH_KEY_OPT) -P $(SSH_PORT) service.sh $(SSH_USER)@$(SSH_HOST):$(REMOTE_GO_DIR)/service.sh
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "chmod +x $(REMOTE_GO_DIR)/service.sh && bash $(REMOTE_GO_DIR)/service.sh migrate"
	@echo "===> [MIGRATE] Hoàn tất migrations trên server!"

seed-remote-master:
	@echo "------------------------------------------------------------------------"
	@echo "===> [SEED] Chạy seeder master data $(if $(table),cho bảng $(table),toàn bộ) trên $(SSH_USER)@$(SSH_HOST)..."
	@echo "------------------------------------------------------------------------"
	@scp $(SSH_KEY_OPT) -P $(SSH_PORT) service.sh $(SSH_USER)@$(SSH_HOST):$(REMOTE_GO_DIR)/service.sh
	@$(SSH_CMD) $(SSH_USER)@$(SSH_HOST) "chmod +x $(REMOTE_GO_DIR)/service.sh && bash $(REMOTE_GO_DIR)/service.sh seed-master $(table)"
	@echo "===> [SEED] Hoàn tất seeder master data trên server!"
