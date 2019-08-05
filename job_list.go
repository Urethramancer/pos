package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdJobList options.
type CmdJobList struct {
	opt.DefaultHelp
}

// Run list
func (cmd *CmdJobList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

