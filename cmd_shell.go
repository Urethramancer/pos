package main

import (
	"errors"

	"github.com/Urethramancer/pos/internal/shell"
	"github.com/Urethramancer/signor/opt"
)

// Cmdshell options.
type CmdShell struct {
	opt.DefaultHelp
}

// Run shell
func (cmd *CmdShell) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	sh, err := shell.New()
	if err != nil {
		return err
	}

	defer sh.Close()
	return sh.Run()
}
