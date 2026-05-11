#!/bin/bash
set -e

echo "➤ Entering backend directory..."
cd backend

echo "➤ Downloading Go dependencies..."
go mod download

echo "➤ Starting the backend server..."
go run main.go
