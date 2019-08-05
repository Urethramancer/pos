package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdJobCancel options.
type CmdJobCancel struct {
	opt.DefaultHelp
}

// Run cancel
func (cmd *CmdJobCancel) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
