# gator

A small CLI that interacts with a Postgres-backed feed database.

## Requirements

- Go (tested with Go 1.26+)
- Postgres (running and reachable from your machine)

## Install the gator CLI

Recommended (install from module path):

```bash
go install github.com/erlanggajuni45/gator@latest
```

This will install the `gator` binary into your `GOBIN` or `$(go env GOPATH)/bin`. Make sure that directory is in your `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
# or, if you set GOBIN:
export PATH="$(go env GOBIN):$PATH"
```

For local development (without installing to your bin), run the program directly from source:

```bash
go run main.go <command> [args...]
```

## Configuration

gator reads a JSON config file from your home directory at `~/.gatorconfig.json`.

Create `~/.gatorconfig.json` with at least the `db_url` field. Example:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gatordb?sslmode=disable",
  "current_user_name": ""
}
```

- `db_url`: the Postgres connection string used by the app.
- `current_user_name`: optional. Commands that require a logged-in user will use this name. You can set it via the `login` command.

The code that reads/writes this file is at `internal/config/index.go`; the filename is `~/.gatorconfig.json`.

## Database setup

The repository uses generated SQL code (sqlc) and expects the following tables. If you don't have migrations, you can create the tables manually. Example schema (basic; adjust types/constraints as needed):

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  name TEXT UNIQUE NOT NULL
);

CREATE TABLE feeds (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  url TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  last_fetched_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE feed_follows (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE,
  UNIQUE(user_id, feed_id)
);

CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  description TEXT,
  published_at TIMESTAMP WITH TIME ZONE,
  feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE
);
```

Adjust `uuid` extension and the `uuid_generate_v4()` usage according to your Postgres setup (you can instead insert UUIDs from the application).

## Run the program

After configuring `~/.gatorconfig.json` and ensuring Postgres is reachable:

- If you installed the binary:

  ```bash
  gator <command> [args...]
  ```

- From source:
  ```bash
  go run main.go <command> [args...]
  ```

If you run without any command the program will exit with usage info:

```
usage: cli <command> [args...]
```

## Common commands

`main.go` registers these commands (available from the CLI):

- `register <username>` — create a new user.
- `login <username>` — mark a user as the current user in the config.
- `reset` — delete all users (use with caution; this calls `ResetUsers`).
- `users` — list users.
- `agg` — run aggregation / fetch (triggers feed fetching; exact behavior is implemented in code).
- `addfeed <name> <url>` — add a feed (requires a logged-in user).
- `feeds` — list feeds (shows feed owner).
- `follow <feed-id|feed-url|feed-name>` — follow a feed (requires logged-in user).
- `following` — list feeds the current user is following.
- `unfollow <feed-id|feed-url|feed-name>` — unfollow a feed (requires logged-in user).
- `browse [limit]` — browse posts for the current user (requires logged-in user).

Notes:

- The exact accepted arguments for `follow` / `unfollow` / `addfeed` depend on the handlers' implementation — these are the commands that exist; check handlers in code if you need precise argument formats.
- `login` updates `current_user_name` in `~/.gatorconfig.json`.

## Where in the code

- CLI entrypoint: `main.go`
- Config read/write: `internal/config/index.go` (file name: `.gatorconfig.json` in home)
- Database access (sqlc generated): `internal/database/*`
- Module path: `github.com/erlanggajuni45/gator` (`go.mod`)

## Troubleshooting

- "could not connect to database": verify `db_url`, Postgres is running, and credentials/host/port are correct.
- If `gator` binary isn't found after `go install`, ensure `$(go env GOPATH)/bin` or `$(go env GOBIN)` is in your `PATH`.

---
