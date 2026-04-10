#!/usr/bin/sh
set -e
migrate \
  -source "file://db/migrations" \
  -database "postgresql://postgres:password@localhost:5432/todo?sslmode=disable" $1
