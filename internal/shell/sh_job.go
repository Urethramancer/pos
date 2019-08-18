package shell

import (
	"strconv"

	"github.com/Urethramancer/pos/internal/database"
)

func (sh *Shell) jobCommands(args []string) {
	if len(args) == 0 {
		sh.m("list\t\tList all jobs, optionally only for a client ID.")
		sh.m("add\t\tAdd job. Takes client ID and currency code as arguments.")
		sh.m("edit <id>\t\tModify a job. Same arguments as adding.")
		sh.m("remove <id>\tRemove a job by ID. Will also remove all related tasks.")
		sh.m("show <id>\tShow details and tasks for a job.")
		return
	}

	cmd := args[0]
	args = args[1:]
	switch cmd {
	case "add":
		if len(args) < 2 {
			sh.m("You need to specify a client ID and currency code. Example:\n\tjob add 1 GBP")
		} else {
			sh.addJob(args)
		}
	case "edit":
	case "remove":
	case "list":
	case "show":
	}
}

func (sh *Shell) addJob(args []string) {
	client, err := strconv.Atoi(args[0])
	if err != nil {
		sh.e("Error converting ID: %s", err.Error())
		return
	}

	j := &database.Job{
		Client:   int64(client),
		Currency: args[1],
	}
	id, err := sh.db.AddJob(j)
	if err != nil {
		sh.e("Error adding client: %s", err.Error())
		return
	}

	sh.m("Added job with ID %d to client %d.", id, client)
}
