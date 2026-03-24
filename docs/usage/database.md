# Database

Database shortcuts read `DB_CONNECTION`, `DB_DATABASE`, `DB_USERNAME`, and `DB_PASSWORD` from the project's `.env` and run the appropriate command inside the service container.

## Commands

| Command | Description |
|---|---|
| `lerd db:create [name]` | Create a database and a `<name>_testing` database for the current project |
| `lerd db:import [-d name] <file.sql>` | Import a SQL dump (defaults to site DB from `.env`) |
| `lerd db:export [-d name] [-o file.sql]` | Export a database to a SQL dump (defaults to site DB from `.env`) |
| `lerd db:shell` | Open an interactive MySQL or PostgreSQL shell for the current project |
| `lerd db create [name]` | Same as `db:create` (subcommand form) |
| `lerd db import [-d name] <file.sql>` | Same as `db:import` (subcommand form) |
| `lerd db export [-d name]` | Same as `db:export` (subcommand form) |
| `lerd db shell` | Same as `db:shell` (subcommand form) |

---

## `lerd db:create` name resolution

Name is resolved in this order (first match wins):

1. Explicit `[name]` argument
2. `DB_DATABASE` from the project's `.env`
3. Project name derived from the registered site name (or directory name)

A `<name>_testing` database is always created alongside the main one. If a database already exists the command reports it instead of failing.

Supports `DB_CONNECTION=mysql` / `mariadb` (via `lerd-mysql`) and `pgsql` / `postgres` (via `lerd-postgres`).
