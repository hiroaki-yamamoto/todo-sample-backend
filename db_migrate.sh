#!/usr/bin/sh
set -e
migrate \
  -source "file://db/migrations" \
  -database "postgresql://postgres:password@localhost:5432/$1?sslmode=disable" \
  $2
