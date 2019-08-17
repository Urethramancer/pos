package main

import (
	"errors"

	"github.com/Urethramancer/pos/internal/shell"
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
	sh, err := shell.New()
	if err != nil {
		return err
	}

	return sh.Setup()
}
