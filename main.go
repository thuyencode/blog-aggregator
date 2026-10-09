package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joho/godotenv"
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

	ctx := context.Background()
	dbString := os.Getenv("DB_STRING")
	conn, err := pgx.Connect(ctx, dbString)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	state := &state{config: &cfg, queries: database.New(conn)}

	cmds := make(commands)
	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)

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

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("a single argument is expected")
	}

	username := cmd.Args[0]
	user, err := s.queries.GetUser(context.Background(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("username %q doesn't exist", username)
		}
		return err
	}

	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}

	log.Print("the username has been set.")
	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("a single argument is expected")
	}

	username := cmd.Args[0]

	user, err := s.queries.CreateUser(context.Background(), username)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return fmt.Errorf("username %q is already registered", username)
			default:
				return pgErr
			}
		}
		return err
	}

	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}

	log.Printf("the username has been created: %+v", user)
	return nil
}
