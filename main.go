package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/uLuKaiDev/Gator/internal/app"
	"github.com/uLuKaiDev/Gator/internal/config"
	"github.com/uLuKaiDev/Gator/internal/database"
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

	dbConn, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
	}

	dbQueries := database.New(dbConn)

	state := &app.State{
		DB:     dbQueries,
		DBConn: dbConn,
		Config: &cfg,
	}

	cmd := app.Command{
		Name: os.Args[1],
		Args: os.Args[2:], // everything after the command name
	}

	cmds := app.NewCommands()
	cmds.Register("login", app.HandlerLogin)
	cmds.Register("register", app.HandlerRegister)
	cmds.Register("db-reset", app.HandlerDBReset)
	cmds.Register("reset", app.HandlerDeleteUsers)
	cmds.Register("users", app.HandlerGetUsers)
	cmds.Register("agg", app.HandlerAgg)
	cmds.Register("feeds", app.HandlerListFeeds)
	cmds.Register("addfeed", app.MiddlewareLoggedIn(app.HandlerAddFeed))
	cmds.Register("follow", app.MiddlewareLoggedIn(app.HandlerFollow))
	cmds.Register("following", app.MiddlewareLoggedIn(app.HandlerFollowing))

	if err := cmds.Run(state, cmd); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
