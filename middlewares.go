package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/thuyencode/blog-aggregator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) commandHandler {
	return func(s *state, cmd command) error {
		user, err := s.queries.GetUserByName(context.Background(), s.config.CurrentUserName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return NonExistUsernameError{s.config.CurrentUserName}
			}
			return err
		}

		return handler(s, cmd, user)
	}
}
