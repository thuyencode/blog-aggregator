package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/thuyencode/blog-aggregator/internal/database"
	"github.com/thuyencode/blog-aggregator/internal/feed"
)

type LackArgsError struct {
	ArgNames []string
	Command  command
}

func (e LackArgsError) Error() string {
	var msg strings.Builder
	fmt.Fprintf(&msg, "the %q command expects ", e.Command.Name)

	for _, arg := range e.ArgNames {
		fmt.Fprintf(&msg, "[%s] ", arg)
	}

	if len(e.ArgNames) <= 1 {
		msg.WriteString("as argument")
	} else {
		msg.WriteString("as arguments")
	}

	return msg.String()
}

type NonExistUsernameError struct {
	Username string
}

func (e NonExistUsernameError) Error() string {
	return fmt.Sprintf("username %q doesn't exist", e.Username)
}

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return LackArgsError{[]string{"username"}, cmd}
	}

	username := cmd.Args[0]
	user, err := s.queries.GetUser(context.Background(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NonExistUsernameError{username}
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
		return LackArgsError{[]string{"username"}, cmd}
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
	return s.queries.DeleteAllUsers(context.Background())
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

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return LackArgsError{[]string{"name", "url"}, cmd}
	}

	username := s.config.CurrentUserName
	user, err := s.queries.GetUser(context.Background(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NonExistUsernameError{username}
		}
		return err
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.queries.CreateFeed(
		context.Background(),
		database.CreateFeedParams{Name: name, Url: url, UserID: user.ID},
	)
	if err != nil {
		if pqErr := pq.As(err, pqerror.UniqueViolation); pqErr != nil {
			return fmt.Errorf("url %q is already registered", url)
		}
		return err
	}

	log.Printf("new feed added:")
	fmt.Printf("%+v", feed)

	return nil
}

func handleGetFeeds(s *state, _ command) error {
	feeds, err := s.queries.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	log.Print("feeds list:")
	for _, feed := range feeds {
		fmt.Printf("* %+v\n", feed)
	}

	return nil
}
