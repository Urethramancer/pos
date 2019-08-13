package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdSetup subcommands.
type CmdSetup struct {
	opt.DefaultHelp
}

// Run client
func (cmd *CmdSetup) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	var err error
	sh, err := NewShell()
	if err != nil {
		return err
	}

	return sh.Setup()
}
