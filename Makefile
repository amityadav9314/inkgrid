# InkGrid Makefile

# Variables
BACKEND_DIR = ./backend
FRONTEND_DIR = ./frontend
BIN_DIR = $(BACKEND_DIR)/bin
BIN_NAME = inkgrid
GO = go
NPM = npm

# Colors for terminal output
RESET = \033[0m
GREEN = \033[32m
YELLOW = \033[33m
BLUE = \033[34m

.PHONY: all clean-backend build-backend run-backend clean-frontend build-frontend run-frontend help

all: build-backend build-frontend

# Backend targets
clean-backend:
	@echo "$(BLUE)Cleaning backend...$(RESET)"
	@rm -rf $(BIN_DIR)
	@echo "$(GREEN)Backend cleaned successfully!$(RESET)"

build-backend: clean-backend
	@echo "$(BLUE)Building backend...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@cd $(BACKEND_DIR) && $(GO) build -o bin/$(BIN_NAME) .
	@echo "$(GREEN)Backend built successfully!$(RESET)"

run-backend:
	@echo "$(BLUE)Running backend...$(RESET)"
	@cd $(BACKEND_DIR) && $(GO) run main.go

# Frontend targets
clean-frontend:
	@echo "$(BLUE)Cleaning frontend...$(RESET)"
	@cd $(FRONTEND_DIR) && rm -rf node_modules build
	@echo "$(GREEN)Frontend cleaned successfully!$(RESET)"

build-frontend:
	@echo "$(BLUE)Building frontend...$(RESET)"
	@cd $(FRONTEND_DIR) && $(NPM) install && $(NPM) run build
	@echo "$(GREEN)Frontend built successfully!$(RESET)"

run-frontend:
	@echo "$(BLUE)Running frontend...$(RESET)"
	@cd $(FRONTEND_DIR) && $(NPM) start

# Docker targets
docker-build:
	@echo "$(BLUE)Building Docker images...$(RESET)"
	docker-compose build
	@echo "$(GREEN)Docker images built successfully!$(RESET)"

docker-up:
	@echo "$(BLUE)Starting Docker containers...$(RESET)"
	docker-compose up -d
	@echo "$(GREEN)Docker containers started successfully!$(RESET)"

docker-down:
	@echo "$(BLUE)Stopping Docker containers...$(RESET)"
	docker-compose down
	@echo "$(GREEN)Docker containers stopped successfully!$(RESET)"

# Development targets
dev: run-backend run-frontend

# Help target
help:
	@echo "$(YELLOW)Available targets:$(RESET)"
	@echo "  all             - Build both backend and frontend"
	@echo "  clean-backend   - Clean backend build artifacts"
	@echo "  build-backend   - Build the backend application"
	@echo "  run-backend     - Run the backend application"
	@echo "  clean-frontend  - Clean frontend build artifacts"
	@echo "  build-frontend  - Build the frontend application"
	@echo "  run-frontend    - Run the frontend application in development mode"
	@echo "  docker-build    - Build Docker images"
	@echo "  docker-up       - Start Docker containers"
	@echo "  docker-down     - Stop Docker containers"
	@echo "  dev             - Run both backend and frontend for development"
	@echo "  help            - Show this help message"
