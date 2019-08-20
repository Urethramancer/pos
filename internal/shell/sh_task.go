package shell

func (sh *Shell) taskCommands(args []string) {
	if len(args) == 0 {
		sh.m("list [<id>]\t\tList all tasks, optionally only for a job ID.")
		sh.m("add <id>\t\tAdd task for a job ID. Details will be prompted for.")
		sh.m("edit <id>\tModify a task. Works just like add.")
		sh.m("addtime <id> <amount>\tAdd time to a task.")
		sh.m("cost <id> <amount>\tChange cost for a task.")
		sh.m("remove <id...>\tRemove tasks by ID.")
		sh.m("show <id>\tShow details for a task.")
		sh.m("start <id>\t\tStart timer for a task.")
		sh.m("stop <id>\t\tStop timer for a task.")
		return
	}

	cmd := args[0]
	args = args[1:]
	switch cmd {
	case "add":
	case "edit":
	case "addtime":
	case "cost":
	case "remove":
	case "list":
	case "show":
	case "start":
	case "stop":
	}
}
