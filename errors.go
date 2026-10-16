package main

import (
	"fmt"
	"strings"
)

type NonExistUsernameError struct {
	Username string
}

func (e NonExistUsernameError) Error() string {
	return fmt.Sprintf("username %q doesn't exist", e.Username)
}

type LackArgsError struct {
	ArgNames []string
	Command  command
}

func (e LackArgsError) Error() string {
	var msg strings.Builder
	fmt.Fprintf(&msg, "the %q command requires ", e.Command.Name)

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
