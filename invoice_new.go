package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdInvoiceNew options.
type CmdInvoiceNew struct {
	opt.DefaultHelp
}

// Run new
func (cmd *CmdInvoiceNew) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
