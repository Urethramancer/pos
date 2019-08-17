package shell

import (
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/Urethramancer/pos/internal/database"
	"github.com/Urethramancer/signor/stringer"
)

func (sh *Shell) clientCommands(args []string) {
	// commands := []string{"list", "add", "remove", "show", "find"}
	if len(args) == 0 {
		sh.m("list\t\tList all clients.")
		sh.m("add\t\tAdd client. You will be asked for details to add.")
		sh.m("edit <id>\t\tModify client. You will be asked for details to change.")
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

	case "edit":
		if len(args) == 0 {
			sh.m("You must specify a client ID to edit.")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			sh.e("Error parsing ID: %s", err.Error())
			break
		}

		sh.editClient(int64(id))

	case "remove":
		if len(args) == 0 {
			sh.m("You must specify a client ID to remove.")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			sh.e("Error parsing ID: %s", err.Error())
			break
		}

		sh.removeClient(int64(id))

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
		sh.findClients(args[0])

	default:
		sh.m("Unknown arguments.")
	}
}

func (sh *Shell) listClients() {
	list, err := sh.db.GetAllClients()
	if err != nil {
		sh.e("Error retrieving client list: %s", err.Error())
		return
	}

	printClients(list)
}

func (sh *Shell) addClient() {
	c := sh.promptClient(nil)
	if c == nil {
		return
	}

	id, err := sh.db.AddClient(c)
	if err != nil {
		sh.e("Error adding client: %s", err.Error())
		return
	}

	sh.m("Added client %s with ID %d.", c.Company, id)
}

func (sh *Shell) editClient(id int64) {
	c := sh.db.GetClient(id)
	if c == nil {
		sh.e("No client with that ID.")
		return
	}

	c = sh.promptClient(c)
	if c == nil {
		return
	}

	err := sh.db.UpdateClient(c)
	if err != nil {
		sh.e("Error updating client: %s", err.Error())
		return
	}

	sh.m("Updated client %s with ID %d.", c.Company, c.ID)
}

func (sh *Shell) promptClient(c *database.Client) *database.Client {
	var err error
	if c == nil {
		c = &database.Client{}
	}

	var s string
	if c.Company == "" {
		s, err = sh.Prompt("Company name: ")
	} else {
		s, err = sh.Prompt("Company name [" + c.Company + "]: ")
	}

	if err != nil {
		return nil
	}

	if s != "" {
		c.Company = s
	}

	if c.Email == "" {
		s, err = sh.Prompt("E-mail: ")
	} else {
		s, err = sh.Prompt("E-mail [" + c.Email + "]: ")
	}

	if err != nil {
		return nil
	}

	if s != "" {
		c.Email = s
	}

	if c.Phone == "" {
		s, err = sh.Prompt("Phone: ")
	} else {
		s, err = sh.Prompt("Phone [" + c.Phone + "]: ")
	}

	if err != nil {
		return nil
	}

	if s != "" {
		c.Phone = s
	}

	if c.Address == "" {
		s, err = sh.Prompt("Address: ")
	} else {
		s, err = sh.Prompt("Address [" + c.Address + "]: ")
	}

	if err != nil {
		return nil
	}

	if s != "" {
		c.Address = s
	}

	if c.VATID == "" {
		s, err = sh.Prompt("VAT registration/ID: ")
	} else {
		s, err = sh.Prompt("VAT registration/ID [" + c.VATID + "]: ")
	}

	if err != nil {
		return nil
	}

	if s != "" {
		c.VATID = s
	}

	return c
}

func (sh *Shell) removeClient(id int64) {
	q := "DELETE FROM public.clients WHERE id=$1;"
	_, err := sh.db.Exec(q, id)
	if err != nil {
		sh.e("Error removing client: %s", err.Error())
		return
	}

	sh.m("Removed client %d.", id)
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

func (sh *Shell) findClients(keyword string) {
	list, err := sh.db.GetClients(keyword)
	if err != nil {
		sh.e("Error retrieving clients: %s", err.Error())
		return
	}

	if len(list) == 0 {
		sh.m("No clients found.")
		return
	}

	printClients(list)
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
