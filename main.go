package main

import (
	"os"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
)

// setup, get, set, invoice, customer, job
// setup: no subcommands
// get, set: all vars
// invoice: new, finish, publish, cancel, list
// customer: new, remove, list
// job: new, cancel, finish, list

// Options holds all the tool commands.
var Options struct {
	opt.DefaultHelp
	Setup   CmdSetup   `command:"setup" help:"Set up the basic configuration."`
	Client  CmdClient  `command:"client" help:"Subcommands for client management." aliases:"cl"`
	Invoice CmdInvoice `command:"invoice" help:"Subcommands for invoice generation and publishing." aliases:"inv"`
	Job     CmdJob     `command:"job" help:"Job/work order management."`
	Shell   CmdShell   `command:"shell" help:"Run interactive shell." aliases:"sh"`
}

// ParseOptions could be renamed to main(), or called from main.
func main() {
	a := opt.Parse(&Options)
	if Options.Help || len(os.Args) < 2 {
		a.Usage()
		return
	}

	var err error
	e := log.Default.Err
	err = a.RunCommand(false)
	if err != nil {
		e("Error running: %s", err.Error())
		os.Exit(2)
	}

}
