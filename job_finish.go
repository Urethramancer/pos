package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdJobFinish options.
type CmdJobFinish struct {
	opt.DefaultHelp
}

// Run finish
func (cmd *CmdJobFinish) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

