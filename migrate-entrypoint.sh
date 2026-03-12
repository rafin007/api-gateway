#!/bin/sh
set -e
export GOOSE_DBSTRING="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=require"
export GOOSE_DRIVER=postgres
export GOOSE_MIGRATION_DIR=/migrations
exec goose up