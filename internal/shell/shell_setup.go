package shell

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Urethramancer/cross"
	"github.com/Urethramancer/pos/internal/database"
	"github.com/Urethramancer/signor/log"
	"github.com/peterh/liner"
)

func (sh *Shell) Setup() error {
	list := []string{
		"Your name",
		"E-mail",
		"Database host",
		"Database port",
		"Database name",
		"Database username",
		"Database password",
		"Company name",
		"Company ID (org. #)",
		"VAT percentage",
		"Invoice prefix",
	}

	results := []string{
		sh.cfg.Name,
		sh.cfg.Email,
		sh.cfg.Host,
		sh.cfg.Port,
		sh.cfg.DBName,
		sh.cfg.Username,
		sh.cfg.Password,
		sh.cfg.Company,
		sh.cfg.CompanyID,
		sh.cfg.VAT,
		sh.cfg.InvoicePrefix,
	}

	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	defer line.Close()

	m := log.Default.Msg
	var x string
	var err error
	for i, p := range list {
		x, err = line.Prompt(fmt.Sprintf("%s [%s]: ", p, results[i]))
		if err != nil {
			if err == liner.ErrPromptAborted {
				return nil
			}

			return err
		}

		if x != "" {
			line.AppendHistory(x)
			results[i] = x
		}
	}

	sh.cfg.Name = results[0]
	sh.cfg.Email = results[1]
	sh.cfg.Host = results[2]
	sh.cfg.Port = results[3]
	sh.cfg.DBName = results[4]
	sh.cfg.Username = results[5]
	sh.cfg.Password = results[6]
	sh.cfg.Company = results[7]
	sh.cfg.CompanyID = results[8]
	sh.cfg.VAT = results[9]
	sh.cfg.InvoicePrefix = results[10]

	x, err = line.Prompt(fmt.Sprintf("First invoice number [%d]", sh.cfg.FirstInvoice))
	if err != nil {
		return err
	}

	if x != "" {
		n, err := strconv.Atoi(x)
		if err != nil {
			return err
		}

		sh.cfg.FirstInvoice = n
	}

	x, err = line.Prompt("Year prefix (yes/no) [no]")
	if err != nil {
		return err
	}

	switch x {
	case "y", "yes", "Y", "YES", "t", "T", "true", "TRUE":
		sh.cfg.YearPrefix = true
	case "":
		break
	default:
		sh.cfg.YearPrefix = false
	}

	m("\nEnter the address to print on invoices. Enter '.' to end input.")
	oldaddr := strings.Split(sh.cfg.Address, "\n")
	addr := []string{}
	loop := true
	i := 1
	for loop {
		old := ""
		if len(oldaddr) >= i {
			old = oldaddr[i-1]
		}
		x, err = line.Prompt(fmt.Sprintf("Address line %d [%s]: ", i, old))
		if err != nil {
			if err == liner.ErrPromptAborted {
				return nil
			} else {
				return nil
			}
		}

		if x == "." {
			break
		}

		if x == "" {
			if old == "" {
				break
			}

			x = old
		}

		addr = append(addr, x)
		i++
	}

	sh.cfg.Address = strings.Join(addr, "\n")
	err = sh.cfg.Save(cross.ConfigName("config.json"))
	if err != nil {
		return err
	}

	err = database.TestDBHost(sh.cfg)
	if err != nil {
		return err
	}

	m("DB server pinged OK. Ensuring database exists.")
	err = database.EnsureDBExists(sh.cfg)
	if err != nil {
		return err
	}

	m("Database OK.")
	return nil
}
