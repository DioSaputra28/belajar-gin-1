#!/bin/sh
set -e

echo "🚀 Starting deployment process..."

# Database connection details from environment variables
DB_URL="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

echo "📦 Installing golang-migrate..."
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

echo "🔄 Running database migrations..."
migrate -path database/migrations -database "$DB_URL" -verbose up

echo "✅ Migrations completed successfully!"

echo "🎯 Starting application..."
exec /bin/app
