# Debezium Outbox Go Example

An example of the Outbox pattern using Debezium, and a Go producer and consumer application.

## Project Dependencies

- Docker
- [Task](https://taskfile.dev/)
- Rest Client (VS Code Extension)

## Migrations

Create a new empty migration
```
migrate create -ext sql -dir db/migrations -seq create_users_table
```