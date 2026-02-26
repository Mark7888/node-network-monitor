# Database Migrations

Migrations are managed with [pressly/goose v3](https://github.com/pressly/goose). SQL files live under `internal/db/<driver>/migrations/` and are embedded into the binary at compile time — no external files are needed at runtime.

Migrations run automatically on startup. Existing databases created before goose was introduced are detected and fast-forwarded without re-executing already-applied SQL.

---

## Adding a migration

1. Add a numbered SQL file to **both** `internal/db/sqlite/migrations/` and `internal/db/postgres/migrations/`:
   ```
   00004_my_change.sql
   ```
2. Use goose annotations — always include a `Down` section:
   ```sql
   -- +goose Up
   ALTER TABLE nodes ADD COLUMN my_column TEXT;

   -- +goose Down
   ALTER TABLE nodes DROP COLUMN my_column;
   ```
3. Increment `latestMigrationVersion` in `sqlite/schema.go` and `postgres/schema.go`.
4. Rebuild — the file is embedded automatically.

> Never edit a migration file that has already been released. Add a new one instead.

---

## Upgrading (Docker)

No Dockerfile or docker-compose changes are needed. A standard update is all that is required:

```bash
docker compose pull && docker compose up -d
```

---

## Downgrading

Use the goose CLI against the live database **before** rolling back the binary.

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest

# SQLite — roll back one version
goose -dir internal/db/sqlite/migrations sqlite3 ./data/db.sqlite down

# PostgreSQL — roll back one version
goose -dir internal/db/postgres/migrations postgres "host=localhost user=postgres dbname=speedtest sslmode=disable" down
```

Then redeploy the older binary or image.
