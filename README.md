# 🐊 Gator RSS CLI

Gator is a command-line RSS reader written in Go. It allows users to register, log in, follow feeds, and browse aggregated posts — all from the terminal.

---

## Features

* User registration & login
* Reset & delete user accounts
* Add & list RSS feeds
* Follow/unfollow feeds
* Browse latest posts from followed feeds
* Background feed aggregation (polls feeds at intervals)

---

## Project Structure

```
~/
├── main.go                  # CLI entrypoint
├── .gitignore
├── go.mod / go.sum          # Go module files
├── readme.md
├── sqlc.yaml                # SQLC config
├── internal/
│   ├── app/                 # Handlers and command registrations
│   ├── config/              # App configuration logic
│   └── database/            # SQLC generated database code
├── sql/
│   ├── queries/             # SQL queries consumed by SQLC
│   └── schema/              # Goose-compatible DB migrations
```

---

## Database Schema Overview

### `users`

* `id`, `created_at`, `updated_at`, `name`, `email`, `password_hash`

### `feeds`

* `id`, `created_at`, `updated_at`, `name`, `url`, `user_id`

### `feed_follows`

* `id`, `created_at`, `updated_at`, `user_id`, `feed_id`

### `posts`

* `id`, `created_at`, `updated_at`, `title`, `url`, `description`, `published_at`, `feed_id`

> Post uniqueness is enforced via `feed_id + url`.

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

> Commands that require a user session are protected by middleware.

---

## Time Display

Timestamps are stored in full ISO8601 format (UTC with timezone info), but in the CLI they're rendered using a custom `humanizeTime()` function that converts time to relative strings:

```text
Feed: Boot.dev Blog
Title: Learn Go with CLI Projects
URL: <url will be here, available with cmd+click>
Published At: 7 days ago
```

---

## Development Takeaways

* Use `sqlc generate` to regenerate Go code after editing SQL queries.
* Place all `.sql` files for queries in `sql/queries/`.
* Use `-- name:` comments in SQL files to define generated method names.
* Keep `posts.published_at` in UTC to allow consistent sorting.
* Feed aggregation errors like duplicates are expected and logged but not fatal.

---

## Additional Information

* Since `agg` blocks the terminal (runs a loop), open a second terminal tab for testing other commands.
* Add more feeds to better test the `browse` command.
* The `browse` query was extended to join the `feeds` table so that each post shows the `feed_name`.

---

## Possible Future Enhancements

* [ ] Add post "read" tracking
* [ ] Feed deletion and editing

---

## Credits

Built by \[uLuKaiDev], using:

* Go
* SQLC
* PostgreSQL
* Goose for migrations
* Boot.dev course 

---