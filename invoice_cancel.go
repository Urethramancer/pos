package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdInvoiceCancel options.
type CmdInvoiceCancel struct {
	opt.DefaultHelp
}

// Run cancel
func (cmd *CmdInvoiceCancel) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
