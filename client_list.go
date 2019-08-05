package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdClientList options.
type CmdClientList struct {
	opt.DefaultHelp
}

// Run list
func (cmd *CmdClientList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

