package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thuyencode/blog-aggregator/internal/config"
	"github.com/thuyencode/blog-aggregator/internal/database"
)

type state struct {
	config  *config.Config
	queries *database.Queries
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dbString := os.Getenv("DB_STRING")
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	state := &state{config: &cfg, queries: database.New(db)}

	cmds := make(commands)
	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleGetUsers)
	cmds.register("agg", handleAggregator)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("need at least 2 args")
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	if err = cmds.run(state, command{Name: cmdName, Args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
}
