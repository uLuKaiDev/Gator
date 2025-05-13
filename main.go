package main

import (
	"fmt"
	"os"

	"github.com/uLuKaiDev/Gator/internal/app"
	"github.com/uLuKaiDev/Gator/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: no command provided")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	state := &app.State{Config: &cfg}

	cmd := app.Command{
		Name: os.Args[1],
		Args: os.Args[2:], // everything after the command name
	}

	cmds := app.NewCommands()
	cmds.Register("login", app.HandlerLogin)

	if err := cmds.Run(state, cmd); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
