#!/bin/bash

chown -R www-data:www-data /app

go mod download

go run cmd/app/main.go
