package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdInvoicePublish options.
type CmdInvoicePublish struct {
	opt.DefaultHelp
}

// Run publish
func (cmd *CmdInvoicePublish) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
