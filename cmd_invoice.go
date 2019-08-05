package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdInvoice subcommands.
type CmdInvoice struct {
	opt.DefaultHelp
	Cancel  CmdInvoiceCancel  `command:"cancel" help:"<command help>" aliases:"rm"`
	Finish  CmdInvoiceFinish  `command:"finish" help:"<command help>" aliases:"done"`
	List    CmdInvoiceList    `command:"list" help:"<command help>" aliases:"ls"`
	New     CmdInvoiceNew     `command:"new" help:"<command help>"`
	Publish CmdInvoicePublish `command:"publish" help:"<command help>" aliases:"pub"`
}

// Run invoice
func (cmd *CmdInvoice) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
