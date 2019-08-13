package main

import (
	"io"
	"os"
	"sort"
	"strings"

	"github.com/Urethramancer/cross"
	"github.com/Urethramancer/signor/log"
	"github.com/peterh/liner"
)

// Shell for interactive management.
type Shell struct {
	cfg     *Config
	history string
	l       *liner.State
	m       func(string, ...interface{})
	e       func(string, ...interface{})
}

// NewShell creates an interactive shell structure and loads the configuration.
func NewShell() (*Shell, error) {
	sh := &Shell{}
	sh.cfg = &Config{
		Name:          "",
		Email:         "",
		Host:          "localhost",
		Port:          "5432",
		DBName:        "invoices",
		Username:      "postgres",
		Password:      "postgres",
		Company:       "",
		CompanyID:     "",
		VAT:           "0",
		InvoicePrefix: "",
	}

	cross.SetConfigPath("pos")
	fn := cross.ConfigName("config.json")
	err := sh.cfg.Load(fn)
	if err != nil {
		log.Default.Msg("Creating new configuration file.")
		err = sh.cfg.Save(fn)
		if err != nil {
			return nil, err
		}
	}

	sh.history = cross.ConfigName("history")
	sh.m = log.Default.Msg
	sh.e = log.Default.Err
	return sh, nil
}

func (sh *Shell) LoadHistory() {
	sh.l = liner.NewLiner()
	sh.l.SetHistoryLimit(100)
	f, err := os.Open(sh.history)
	if err != nil {
		return
	} else {
		defer f.Close()
		n, err := sh.l.ReadHistory(f)
		if err != nil {
			sh.e("Couldn't load history: %s", err.Error())
			return
		} else {
			sh.m("Loaded %d history items.", n)
		}
	}
}

// SaveHistory to file.
func (sh *Shell) SaveHistory() {
	f, err := os.OpenFile(sh.history, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		sh.e("Couldn't open history file %s: %s", sh.history, err.Error())
	}

	defer f.Close()
	_, err = sh.l.WriteHistory(f)
	if err != nil {
		sh.e("Couldn't write history: %s", err.Error())
	}
}

// Run the shell until exited or some fatal error happens.
func (sh *Shell) Run() error {
	sh.LoadHistory()
	defer sh.SaveHistory()
	sh.m("Ctrl-D to quit, Tab for command completion, ? for a list of commands, ?? for a full shortcut list.")

	commands := []string{
		"setup",
		"clients",
		"invoices",
		"jobs",
		"tasks",
	}
	sort.Strings(commands)
	sh.l.SetCompleter(func(line string) (c []string) {
		for _, n := range commands {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	defer sh.l.Close()
	running := true
	for running {
		x, err := sh.l.Prompt("> ")
		if err != nil {
			if err == io.EOF {
				sh.m("")
				return nil
			}

			return err
		}

		sh.l.AppendHistory(x)
		switch x {
		case "?":
			sh.m("%s", Tablify(commands, 5))
		case "??":
			showShortcuts()
		case "setup":
			err = sh.Setup()
			sh.m("%s", err)
			if err != nil {
				if err == io.EOF {
					sh.m("Setup aborted.")
				} else {
					sh.e("Error during setup: %s", err.Error())
				}
			}
		}
	}

	return nil
}
