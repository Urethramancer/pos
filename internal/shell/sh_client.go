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
		sh.m("remove <id>\tRemove client by ID(s).")
		sh.m("show <id>\tShow details for specified client ID.")
		sh.m("find <keyword>\tSearch for clients matching the text.")
		return
	}

	cmd := args[0]
	args = args[1:]
	switch cmd {
	case "add":
		sh.addClient()

	case "edit":
		if len(args) == 0 {
			sh.m("You must specify a client ID to edit.")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			sh.e("%s: %s", ErrParseID, err.Error())
			break
		}

		sh.editClient(int64(id))

	case "remove":
		if len(args) == 0 {
			sh.m("You must specify one or more client IDs to remove.")
			return
		}

		for _, x := range args {
			id, err := strconv.Atoi(x)
			if err != nil {
				sh.e("%s: %s", ErrParseID, err.Error())
				break
			}
			sh.removeClient(int64(id))
		}

	case "list":
		sh.listClients()

	case "show":
		if len(args) == 0 {
			sh.m("You must specify one or more client IDs to show.")
			return
		}
		sh.showClients(args)

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

func (sh *Shell) addClient() {
	c := sh.promptClient(nil)
	if c == nil {
		return
	}

	id, err := sh.db.AddClient(c)
	if err != nil {
		sh.e("%s: %s", ErrAddClient, err.Error())
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
		sh.e("%s: %s", ErrUpdateClient, err.Error())
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
	s, err = sh.strPrompt(strCompanyName, c.Company)
	if err != nil {
		return nil
	}

	if s != "" {
		c.Company = s
	}

	s, err = sh.strPrompt(strEmail, c.Email)
	if err != nil {
		return nil
	}

	if s != "" {
		c.Email = s
	}

	s, err = sh.strPrompt(strPhone, c.Phone)
	if err != nil {
		return nil
	}

	if s != "" {
		c.Phone = s
	}

	s, err = sh.strPrompt(strAddress, c.Address)
	if err != nil {
		return nil
	}

	if s != "" {
		c.Address = s
	}

	s, err = sh.strPrompt(strVATID, c.VATID)
	if err != nil {
		return nil
	}

	if s != "" {
		c.VATID = s
	}

	return c
}

func (sh *Shell) removeClient(id int64) {
	err := sh.db.RemoveClient(id)
	if err != nil {
		sh.e("%s: %s", ErrRemoveClient, err.Error())
		return
	}

	sh.m("Removed client %d.", id)
}

func (sh *Shell) listClients() {
	list, err := sh.db.GetAllClients()
	if err != nil {
		sh.e("%s: %s", ErrGetClientList, err.Error())
		return
	}

	sh.printClients(list)
}

func (sh *Shell) showClients(idlist []string) {
	var list []*database.Client
	for _, id := range idlist {
		x, err := strconv.Atoi(id)
		if err != nil {
			return
		}
		c := sh.db.GetClient(int64(x))
		if c == nil {
			sh.m("No client with ID %s.", id)
		} else {
			list = append(list, c)
		}
	}
	if len(list) == 0 {
		sh.m("No clients found with supplied ID(s).")
	} else {
		sh.printClients(list)
	}
}

func (sh *Shell) findClients(keyword string) {
	list, err := sh.db.GetClients(keyword)
	if err != nil {
		sh.e("%s: %s", ErrGetClientList, err.Error())
		return
	}

	if len(list) == 0 {
		sh.m("No clients found.")
		return
	}

	sh.printClients(list)
}

func (sh *Shell) printClients(list []*database.Client) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	s := stringer.New()
	s.WriteString("ID\tClient\tE-mail\tPhone\tAddress\tVAT ID\tCreated\n")
	for _, c := range list {
		_, err := s.WriteI(
			c.ID,
			"\t", c.Company,
			"\t", c.Email,
			"\t", c.Phone,
			"\t", c.Address,
			"\t", c.VATID,
			"\t", c.Created.String(),
			"\n",
		)

		if err != nil {
			sh.e("Error printing to stdout: %s", err.Error())
			return
		}
	}
	_, _ = tw.Write([]byte(s.String()))
	tw.Flush()
}
