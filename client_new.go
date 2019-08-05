package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdClientNew options.
type CmdClientNew struct {
	opt.DefaultHelp
}

// Run new
func (cmd *CmdClientNew) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

