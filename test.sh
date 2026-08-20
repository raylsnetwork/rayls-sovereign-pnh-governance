#!/bin/bash

echo "Running tests with coverage..."
go test -coverprofile=coverage.out -coverpkg=./... ./cmd/api/integration

echo "Generating coverage report..."

go tool cover -html=coverage.out -o coverage.html
explorer.exe coverage.html
