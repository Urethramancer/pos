package shell

import (
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
		sh.m("show <id>\tShow details for specified contact ID.")
		sh.m("find <keyword>\tSearch for contacts matching the text.")
		return
	}

	cmd := args[0]
	args = args[1:]
	switch cmd {
	case "add":
		sh.addContact()

	case "edit":
		if len(args) == 0 {
			sh.m("You must specify a contact ID to edit.")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			sh.e("%s: %s", ErrParseID, err.Error())
			break
		}

		sh.editContact(int64(id))

	case "remove":
		if len(args) == 0 {
			sh.m("You must specify one or more contact IDs to remove.")
			return
		}

		for _, x := range args {
			id, err := strconv.Atoi(x)
			if err != nil {
				sh.e("%s: %s", ErrParseID, err.Error())
				break
			}
			sh.removeContact(int64(id))
		}

	case "list":
		sh.listContacts()

	case "show":
		if len(args) == 0 {
			sh.m("You must specify one or more contact IDs to show.")
			return
		}
		sh.showContacts(args)

	case "find":
		if len(args) == 0 {
			sh.m("You must specify a keyword to search for.")
			return
		}
		sh.findContacts(args[0])
	}
}

func (sh *Shell) addContact() {
	c := sh.promptContact(nil)
	if c == nil {
		return
	}

	id, err := sh.db.AddContact(c)
	if err != nil {
		sh.e("%s: %s", ErrAddContact, err.Error())
		return
	}

	sh.m("Added contact %s with ID %d.", c.Name, id)
}

func (sh *Shell) editContact(id int64) {
	c := sh.db.GetContact(id)
	if c == nil {
		sh.e("No contact with that ID.")
		return
	}

	c = sh.promptContact(c)
	if c == nil {
		return
	}

	err := sh.db.UpdateContact(c)
	if err != nil {
		sh.e("%s: %s", ErrUpdateContact, err.Error())
		return
	}

	sh.m("Updated contact %s with ID %d.", c.Name, c.ID)
}

func (sh *Shell) promptContact(c *database.Contact) *database.Contact {
	var err error
	if c == nil {
		c = &database.Contact{}
	}

	var s string
	s, err = sh.strPrompt(strName, c.Name)
	if err != nil {
		return nil
	}

	if s != "" {
		c.Name = s
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

	s, err = sh.intPrompt(strCompanyID, c.Client)
	if err != nil {
		return nil
	}

	if s != "" {
		id, err := strconv.Atoi(s)
		if err != nil {
			sh.e("%s.", ErrConvertID)
			return nil
		}

		client := sh.db.GetClient(int64(id))
		if client == nil {
			sh.e("No company found with ID %d.", id)
			return nil
		}

		c.Client = int64(id)
	} else {
		if c.Client == 0 {
			sh.e("You must specify a company ID.")
			return nil
		}
	}

	return c
}

func (sh *Shell) removeContact(id int64) {
	err := sh.db.RemoveContact(id)
	if err != nil {
		sh.e("%s: %s", ErrRemoveContact, err.Error())
		return
	}

	sh.m("Removed contact %d.", id)
}

func (sh *Shell) listContacts() {
	list, err := sh.db.GetAllContacts()
	if err != nil {
		sh.e("%s: %s", ErrGetContactList, err.Error())
		return
	}

	sh.printContacts(list)
}

func (sh *Shell) showContacts(idlist []string) {
	var list []*database.Contact
	for _, id := range idlist {
		x, err := strconv.Atoi(id)
		if err != nil {
			return
		}
		c := sh.db.GetContact(int64(x))
		if c == nil {
			sh.m("No contact with ID %s.", id)
		} else {
			list = append(list, c)
		}
	}
	if len(list) == 0 {
		sh.m("No clients found with supplied ID(s).")
	} else {
		sh.printContacts(list)
	}
}

func (sh *Shell) findContacts(keyword string) {
	list, err := sh.db.GetContacts(keyword)
	if err != nil {
		sh.e("%s: %s", ErrGetContactList, err.Error())
		return
	}

	if len(list) == 0 {
		sh.m("No clients found.")
		return
	}

	sh.printContacts(list)
}

func (sh *Shell) printContacts(list []*database.Contact) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	s := stringer.New()
	s.WriteString("ID\tClient\tE-mail\tPhone\tClient\tCreated\n")
	for _, c := range list {
		company := "Unknown"
		client := sh.db.GetClient(c.Client)
		if client != nil {
			company = client.Company
		}
		_, err := s.WriteI(
			c.ID,
			"\t", c.Name,
			"\t", c.Email,
			"\t", c.Phone,
			"\t", c.Client, " (", company, ")",
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
