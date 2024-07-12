package main

//go:generate tern migrate --migrations internal/pg/migrations --config internal/pg/migrations/tern.conf
//go:generate sqlc generate -f internal/pg/sqlc.yaml