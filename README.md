# 🐊 Gator RSS CLI

Gator is a command-line RSS reader written in Go. It allows users to register, log in, follow feeds, and browse aggregated posts — all from the terminal.

---

## Prerequisites

Before using Gator, you’ll need:

* **Go** (1.20 or later): [Install Go](https://golang.org/dl/)
* **PostgreSQL**: You’ll need a running Postgres database (e.g. via Docker or local install)
* **`goose`**: Used for running database migrations — install with:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

## Installation

To install the `gator` CLI globally:

```sh
go install github.com/yourusername/gator@latest
```

Make sure your `$GOPATH/bin` is in your `$PATH` so you can run `gator` from anywhere.

---

## Configuration

Before running the app, create a config file at:

```
./config/config.json
```

Example contents:

```json
{
  "db_url": "postgres://postgres:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "ulukai"
}
```

Then run database migrations:

```sh
goose -dir ./sql/schema postgres "your-db-url" up
```

Replace `"your-db-url"` with your actual Postgres connection string.

---

## Running the Program

For development:

```sh
go run .
```

For production (after `go install`):

```sh
gator browse
```

Note: The `agg` command blocks the terminal as it polls feeds — open a second tab to test other commands.

---

## CLI Commands

| Command      | Description                                  |
| ------------ | -------------------------------------------- |
| `register`   | Register a new user                          |
| `login`      | Log in as a user and persist session         |
| `reset`      | Deletes all users                            |
| `db-reset`   | Drops & recreates the database schema        |
| `users`      | Lists all registered users                   |
| `agg 10s`    | Starts scraping feeds every 10s (blocks CLI) |
| `addfeed`    | Add a new RSS feed                           |
| `feeds`      | Lists all feeds in the system                |
| `follow`     | Follow a feed (must be logged in)            |
| `following`  | Lists all feeds you're following             |
| `unfollow`   | Unfollow a feed                              |
| `browse [n]` | Browse most recent `n` posts (default 10)    |

---

## Output Example

```text
Feed: Boot.dev Blog
Title: Learn Go with CLI Projects
URL: https://example.com/post
Published At: 7 days ago
```

> Time is rendered using a custom `humanizeTime()` function.

---

## Project Structure

```
~/
├── main.go                  # CLI entrypoint
├── .gitignore
├── go.mod / go.sum
├── readme.md
├── sqlc.yaml
├── internal/
│   ├── app/                 # Handlers and command registrations
│   ├── config/              # App configuration logic
│   └── database/            # SQLC generated code
├── sql/
│   ├── queries/             # SQL queries used by SQLC
│   └── schema/              # DB migrations (Goose)
```

---

## Tips & Notes

* Use `sqlc generate` after editing SQL files in `sql/queries/`
* `published_at` timestamps are stored in full UTC format
* Duplicate post insertions (same feed\_id + url) are ignored

---

## Future Enhancements / Will-never-look-at-this-again 

* [ ] Track post read/unread state
* [ ] Feed management (rename/delete)
* [ ] Unit tests for core logic

---

## Credits

Built by \[uLuKaiDev], powered by:

* Go
* PostgreSQL
* SQLC
* Goose
* Boot.dev Go course
