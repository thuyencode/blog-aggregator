package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/lib/pq/pqerror"
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
		if errors.Is(err, sql.ErrNoRows) {
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
		if pqErr := pq.As(err, pqerror.UniqueViolation); pqErr != nil {
			return fmt.Errorf("username %q is already registered", username)
		}
		return err
	}

	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}

	log.Printf("the username has been created: %+v", user)
	return nil
}

func handleReset(s *state, _ command) error {
	return s.queries.PruneUsers(context.Background())
}

func handleGetUsers(s *state, _ command) error {
	users, err := s.queries.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.config.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}
