PY ?= python
DJANGO ?= $(PY) manage.py
GO ?= go
EXTRACTOR ?= cd go_modules/data_extractor && $(GO) run cmd/main.go
PORT ?= 8000

.PHONY: db-up run

db-up:
	docker compose up -d --wait
	$(MAKE) db-migrate
	$(EXTRACTOR)

db-down:
	docker compose down -v

db-reset:
	docker compose down -v
	docker compose up postgres --build -d --wait
	$(MAKE) db-migrate

db-migrate:
	$(DJANGO) makemigrations
	$(DJANGO) migrate

run:
	# $(MAKE) db-up
	$(DJANGO) runserver $(PORT)