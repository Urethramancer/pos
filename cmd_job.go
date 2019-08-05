package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdJob subcommands.
type CmdJob struct {
	opt.DefaultHelp
	Cancel	CmdJobCancel	`command:"cancel" help:"<command help>"`
	Finish	CmdJobFinish	`command:"finish" help:"<command help>"`
	List	CmdJobList	`command:"list" help:"<command help>"`
	New	CmdJobNew	`command:"new" help:"<command help>"`
}

// Run job
func (cmd *CmdJob) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

