package shell

import (
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/Urethramancer/signor/stringer"

	"github.com/Urethramancer/pos/internal/database"
)

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
		sh.listClients()

	case "add":
		sh.addClient()

	case "remove":
		if len(args) == 0 {
			sh.m("You must specify a client ID to remove.")
			return
		}
		sh.removeClient()

	case "show":
		if len(args) == 0 {
			sh.m("You must specify a client ID to show.")
			return
		}
		sh.showClient(args[0])

	case "find":
		if len(args) == 0 {
			sh.m("You must specify a keyword to search for.")
			return
		}
		sh.findClients()

	default:
		sh.m("Unknown arguments.")
	}
}

func (sh *Shell) listClients() {

}

func (sh *Shell) addClient() {
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

func (sh *Shell) removeClient() {

}

func (sh *Shell) showClient(id string) {
	x, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	c := sh.db.GetClient(int64(x))
	if c == nil {
		sh.m("No client with that ID.")
	} else {
		printClient(c)
	}
}

func (sh *Shell) findClients() {

}

func printClient(c *database.Client) {
	printClients([]*database.Client{c})
}

func printClients(list []*database.Client) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	s := stringer.New()
	s.WriteString("ID\tClient\tE-mail\tPhone\tAddress\tVAT ID\tCreated\n")
	for _, c := range list {
		s.WriteI(
			c.ID,
			"\t", c.Company,
			"\t", c.Email,
			"\t", c.Phone,
			"\t", c.Address,
			"\t", c.VATID,
			"\t", c.Created.String(),
			"\n",
		)
	}
	_, _ = tw.Write([]byte(s.String()))
	tw.Flush()
}
