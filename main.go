package main

import (
	"errors"
	"log"
	"os"

	"internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	state := &state{config: &cfg}

	cmds := make(commands)
	cmds.register("login", handleLogin)

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
		return errors.New("a single argument is expected: username")
	}

	username := cmd.Args[0]
	if err := s.config.SetUser(username); err != nil {
		return err
	}

	log.Print("the username has been set.")
	return nil
}
