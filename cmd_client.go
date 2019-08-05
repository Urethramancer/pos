package main

import (
	"errors"
	"os"

	"github.com/Urethramancer/signor/opt"
)

// CmdClient subcommands.
type CmdClient struct {
	opt.DefaultHelp
	List   CmdClientList   `command:"list" help:"<command help>" aliases:"ls"`
	New    CmdClientNew    `command:"new" help:"<command help>"`
	Remove CmdClientRemove `command:"remove" help:"<command help>" aliases:"rm,del"`
}

// Run client
func (cmd *CmdClient) Run(in []string) error {
	if cmd.Help || len(os.Args) < 3 {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
