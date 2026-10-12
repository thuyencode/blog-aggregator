package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/thuyencode/blog-aggregator/internal/feed"
)

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

	log.Print("users list:")

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.config.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func handleAggregator(_ *state, _ command) error {
	feed, err := feed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	log.Printf("rss feed:")
	fmt.Printf("%+v", *feed)
	return nil
}
