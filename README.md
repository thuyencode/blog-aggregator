# Blog Aggregator

This is a guided project I did with Boot.dev.

Note: This project doesn't use any form of AI/LLM assistance.

## Prerequisites

You will need to have Go in PostgreSQL installed in order to run this project.
Having Golangci-lint is a good bonus too if you want to use all Makefile targets and git hooks in [`scripts/githooks`](./scripts/githooks).

### Go 1.27.0 or above

This is a Go project, you should have Go installed already. Make sure the minium version is 1.27.0 as declared in [`go.mod`](./go.mod).

Either using your OS package manager, some devtool version managers like ASDF, or following [the official instructions](https://go.dev/doc/install).

### PostgreSQL

It's doesn't matter if your PostgreSQL instance running on the cloud or local, as long as the values in `.env` are correct.

The easiest way to set up PostgreSQL for this project is to use Docker Compose [which I already setup up](./docker-compose.yml). Simply load the enviroment variables in `.env` file to your shell first (in my case I use [loadenv.fish](https://github.com/berk-karaal/loadenv.fish) since I use Fish shell), then run the following command:

```sh
docker compose up -d
```

I'm sure you already know how to install Docker.

### Database migration

You will need [Goose](https://github.com/pressly/goose) to migration the database. Hope I don't need to show you how to install a 🪿

The way you want to run `goose` command depends on you. In my case, I load the enviroment variables from `.env` to shell and then run the following command:

```sh
goose postgres "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" -dir sql/schema up
```

## Usage

Here is the list of all the current commands and their arguments:

- `login`: Login to an account
  - `username`: The username of a user
- `register`: Create a new account
  - `username`: The username of a user
- `reset`: This is for testing only. It'll prune all data
- `users`: List all users
- `agg`: Fetch the RSS feeds
  - `time_between_reqs`: The timeout between each request. The proper formats are [listed here](https://pkg.go.dev/time#ParseDuration).
- `addfeed`: Add a new RSS feed
  - `name`: The name of the RSS feed
  - `url`: The RSS url of the feed
- `feeds`: List all RSS feeds
- `follow`: Follow a RSS feed
  - `url`: The RSS url of the feed
- `following`: List all RSS feeds the current account follows
- `unfollow`: Unfollow a RSS feed
  - `url`: The RSS url of the feed
- `browse`: Show the list of most recent posts
  - `limit`: The number of posts to be shown. Default is 2

The proper command is:

```sh
gator [command] [arg_1] [arg_2]
```
