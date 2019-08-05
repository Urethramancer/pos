package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdInvoiceFinish options.
type CmdInvoiceFinish struct {
	opt.DefaultHelp
}

// Run finish
func (cmd *CmdInvoiceFinish) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
