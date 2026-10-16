package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commandHandler func(s *state, cmd command) error

type commands map[string]commandHandler

func (c commands) run(s *state, cmd command) error {
	cb, ok := c[cmd.Name]
	if !ok {
		return fmt.Errorf("command %q not found", cmd.Name)
	}
	return cb(s, cmd)
}

func (c commands) register(name string, f commandHandler) {
	c[name] = f
}
