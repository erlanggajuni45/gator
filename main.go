package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/erlanggajuni45/gator/internal/config"
	"github.com/erlanggajuni45/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{cfg: &cfg, db: dbQueries}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerResetUsers)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handleFetchFeed)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleListFeeds)
	cmds.register("follow", middlewareLoggedIn(feedFollowsHandler))
	cmds.register("following", middlewareLoggedIn(followingHandler))
	cmds.register("unfollow", middlewareLoggedIn(deleteFeedFollowHandler))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
