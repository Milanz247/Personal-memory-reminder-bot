.PHONY: build run test clean install-service uninstall-service docker-build docker-run help

# Build variables
BINARY_NAME=memory-bot
DOCKER_IMAGE=memory-bot:latest

help: ## Show this help message
	@echo "Memory Storage Bot - Makefile Commands"
	@echo "======================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@CGO_ENABLED=1 go build -tags "fts5" -o $(BINARY_NAME) -ldflags="-s -w" .
	@echo "✅ Build complete!"

run: ## Run the application
	@echo "Starting $(BINARY_NAME)..."
	@./$(BINARY_NAME)

dev: ## Run in development mode
	@echo "Starting in development mode..."
	@go run main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f memories.db*
	@echo "✅ Clean complete!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed!"

install-service: build ## Install systemd service
	@echo "Installing systemd service..."
	@sudo cp memory-bot.service /etc/systemd/system/
	@sudo systemctl daemon-reload
	@sudo systemctl enable memory-bot
	@sudo systemctl start memory-bot
	@echo "✅ Service installed and started!"
	@sudo systemctl status memory-bot

uninstall-service: ## Uninstall systemd service
	@echo "Uninstalling systemd service..."
	@sudo systemctl stop memory-bot
	@sudo systemctl disable memory-bot
	@sudo rm /etc/systemd/system/memory-bot.service
	@sudo systemctl daemon-reload
	@echo "✅ Service uninstalled!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "✅ Docker image built!"

docker-run: ## Run in Docker
	@echo "Running in Docker..."
	@docker run -d \
		--name memory-bot \
		-v $(PWD)/data:/app/data \
		--env-file .env \
		$(DOCKER_IMAGE)
	@echo "✅ Container started!"
	@docker logs -f memory-bot

docker-stop: ## Stop Docker container
	@docker stop memory-bot
	@docker rm memory-bot
	@echo "✅ Container stopped!"

logs: ## View application logs (systemd)
	@sudo journalctl -u memory-bot -f

setup: deps ## Initial setup
	@echo "Setting up project..."
	@chmod +x setup.sh
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "⚠️  Please edit .env and add your TELEGRAM_BOT_TOKEN"; \
	fi
	@echo "✅ Setup complete!"

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

check: format test ## Format code and run tests

backup-db: ## Backup database
	@echo "Backing up database..."
	@cp memories.db memories.db.backup.$(shell date +%Y%m%d_%H%M%S)
	@echo "✅ Database backed up!"

restore-db: ## Restore database from latest backup
	@echo "Restoring database..."
	@cp $(shell ls -t memories.db.backup.* | head -1) memories.db
	@echo "✅ Database restored!"

all: clean deps build ## Clean, download deps, and build
