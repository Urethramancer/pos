package shell

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/Urethramancer/pos/internal/database"
	"github.com/Urethramancer/signor/stringer"
)

func (sh *Shell) contactCommands(args []string) {
	if len(args) == 0 {
		sh.m("list\t\tList all contacts.")
		sh.m("add\t\tAdd contact. You will be asked for details to add.")
		sh.m("edit <id>\t\tModify contact. You will be asked for details to change.")
		sh.m("remove <id>\tRemove contact by ID.")
		sh.m("show <id>\tShow details for specifient contact ID.")
		sh.m("find <keyword>\tSearch for contacts matching the text.")
		return
	}

	cmd := args[0]
	args = args[1:]
	switch cmd {
	case "list":
	case "add":
		sh.addContact()

	case "edit":
	case "remove":
	case "show":
	case "find":
	}
}

func (sh *Shell) addContact() {
	c := sh.promptContact(nil)
	if c == nil {
		return
	}

	id, err := sh.db.AddContact(c)
	if err != nil {
		sh.e("Error adding client: %s", err.Error())
		return
	}

	sh.m("Added client %s with ID %d.", c.Name, id)
}

func (sh *Shell) promptContact(c *database.Contact) *database.Contact {
	var err error
	if c == nil {
		c = &database.Contact{}
	}

	var s string
	if c.Name == "" {
		s, err = sh.Prompt("Name: ")
	} else {
		s, err = sh.Prompt("Name [" + c.Name + "]: ")
	}

	if err != nil {
		return nil
	}

	if s != "" {
		c.Name = s
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

	if c.Client == 0 {
		s, err = sh.Prompt("Company ID: ")
	} else {
		x := fmt.Sprintf("Company ID [%d]: ", c.Client)
		s, err = sh.Prompt(x)
	}

	if err != nil {
		return nil
	}

	if s != "" {
		id, err := strconv.Atoi(s)
		if err != nil {
			sh.e("Error converting ID.")
			return nil
		}

		client := sh.db.GetClient(int64(id))
		if client == nil {
			sh.e("No company found with ID %d.", id)
			return nil
		}

		c.Client = int64(id)
	} else {
		sh.e("You must specify a company ID.")
		return nil
	}

	return c
}

func printContacts(list []*database.Contact) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	s := stringer.New()
	s.WriteString("ID\tClient\tE-mail\tPhone\tAddress\tVAT ID\tCreated\n")
	for _, c := range list {
		s.WriteI(
			c.ID,
			"\t", c.Name,
			"\t", c.Email,
			"\t", c.Phone,
			"\t", c.Client,
			"\t", c.Created.String(),
			"\n",
		)
	}
	_, _ = tw.Write([]byte(s.String()))
	tw.Flush()
}

func (sh *Shell) listContacts() {
	list, err := sh.db.GetAllContacts()
	if err != nil {
		sh.e("Error retrieving client list: %s", err.Error())
		return
	}

	printContacts(list)
}
