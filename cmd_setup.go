package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Urethramancer/cross"
	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
	"github.com/peterh/liner"
)

// CmdSetup subcommands.
type CmdSetup struct {
	opt.DefaultHelp
}

// Run client
func (cmd *CmdSetup) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
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
	}

	cfg := &Config{
		Name:      "",
		Email:     "",
		Host:      "localhost",
		Port:      "5432",
		DBName:    "invoices",
		Username:  "postgres",
		Password:  "postgres",
		Company:   "",
		CompanyID: "",
		VAT:       "0",
	}
	m := log.Default.Msg

	var err error
	cross.SetConfigPath("pos")
	err = cfg.Load(cross.ConfigName("config.json"))
	if err != nil {
		m("Creating new configuration file.")
	}

	results := []string{
		cfg.Name,
		cfg.Email,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Username,
		cfg.Password,
		cfg.Company,
		cfg.CompanyID,
		cfg.VAT,
	}

	var x string
	for i, p := range list {
		x, err = line.Prompt(fmt.Sprintf("%s [%s]: ", p, results[i]))
		if err != nil {
			if err == liner.ErrPromptAborted {
				os.Exit(0)
				return nil
			}
			os.Exit(2)
		}

		if x != "" {
			line.AppendHistory(x)
			results[i] = x
		}
	}

	cfg.Name = results[0]
	cfg.Email = results[1]
	cfg.Host = results[2]
	cfg.Port = results[3]
	cfg.DBName = results[4]
	cfg.Username = results[5]
	cfg.Password = results[6]
	cfg.Company = results[7]
	cfg.CompanyID = results[8]

	m("\nEnter the address to print on invoices. Enter '.' to end input.")
	oldaddr := strings.Split(cfg.Address, "\n")
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
				os.Exit(0)
			} else {
				os.Exit(2)
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

	cfg.Address = strings.Join(addr, "\n")
	err = cfg.Save(cross.ConfigName("config.json"))
	if err != nil {
		return err
	}

	err = testDBHost(cfg)
	if err != nil {
		return err
	}

	m("DB server pinged OK. Ensuring database exists.")
	err = ensureDBExists(cfg)
	if err != nil {
		return err
	}

	return nil
}
