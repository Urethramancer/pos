package shell

func (sh *Shell) clientCommands(args []string) {
	// commands := []string{"list", "add", "remove", "show", "find"}
	if len(args) == 0 {
		sh.m("list\t\tList all clients.")
		sh.m("add\t\tAdd client. You will be asked for details to add.")
		sh.m("remove <id>\tRemove client by ID.")
		sh.m("show <id>\tShow details for specifient client ID.")
		sh.m("find <keyword>\tSearch for clients matching the text.")
		return
	}

	cmd := args[0]
	args = args[1:]
	switch cmd {
	case "list":
		sh.clientList()

	case "add":
		sh.clientAdd()

	case "remove":
		if len(args) == 0 {
			sh.m("You must specify a client ID to remove.")
			return
		}
		sh.clientRemove()

	case "show":
		if len(args) == 0 {
			sh.m("You must specify a client ID to show.")
			return
		}
		sh.clientShow()

	case "find":
		if len(args) == 0 {
			sh.m("You must specify a keyword to search for.")
			return
		}
		sh.clientFind()

	default:
		sh.m("Unknown arguments.")
	}
}

func (sh *Shell) clientList() {

}

func (sh *Shell) clientAdd() {
	name, err := sh.Prompt("Company name: ")
	if err != nil {
		sh.e("Error: %s", err.Error())
		return
	}

	email, err := sh.Prompt("E-mail: ")
	if err != nil {
		sh.e("Error: %s", err.Error())
		return
	}

	phone, err := sh.Prompt("Phone: ")
	if err != nil {
		sh.e("Error: %s", err.Error())
		return
	}

	addr, err := sh.Prompt("Address: ")
	if err != nil {
		sh.e("Error: %s", err.Error())
		return
	}

	vat, err := sh.Prompt("VAT registration/ID: ")
	if err != nil {
		sh.e("Error: %s", err.Error())
		return
	}

	id, err := sh.db.AddClient(name, email, phone, addr, vat)
	if err != nil {
		sh.e("Error adding client: %s", err.Error())
		return
	}

	sh.m("Added client %s with ID %d.", name, id)
}

func (sh *Shell) clientRemove() {

}

func (sh *Shell) clientShow() {

}

func (sh *Shell) clientFind() {

}
