package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdClientRemove options.
type CmdClientRemove struct {
	opt.DefaultHelp
}

// Run remove
func (cmd *CmdClientRemove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

