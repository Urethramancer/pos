package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdJobNew options.
type CmdJobNew struct {
	opt.DefaultHelp
}

// Run new
func (cmd *CmdJobNew) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

