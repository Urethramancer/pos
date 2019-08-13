package main

import (
	"errors"

	"github.com/Urethramancer/signor/log"
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

	sh, err := NewShell()
	if err != nil {
		return err
	}

	return sh.Run()
}

func showShortcuts() {
	log.Default.Msg("Ctrl-A, Home			Move cursor to beginning of line\n" +
		"Ctrl-E, End			Move cursor to end of line\n" +
		"Ctrl-B, Left			Move cursor one character left\n" +
		"Ctrl-F, Right			Move cursor one character right\n" +
		"Ctrl-Left, Alt-B		Move cursor to previous word\n" +
		"Ctrl-Right, Alt-F		Move cursor to next word\n" +
		"Ctrl-D, Del			(if line is not empty) Delete character under cursor\n" +
		"Ctrl-D				(if line is empty) End of File - usually quits application\n" +
		"Ctrl-C				Reset input (create new empty prompt)\n" +
		"Ctrl-L				Clear screen (line is unmodified)\n" +
		"Ctrl-T				Transpose previous character with current character\n" +
		"Ctrl-H, BackSpace		Delete character before cursor\n" +
		"Ctrl-W, Alt-BackSpace		Delete word leading up to cursor\n" +
		"Alt-D				Delete word following cursor\n" +
		"Ctrl-K				Delete from cursor to end of line\n" +
		"Ctrl-U				Delete from start of line to cursor\n" +
		"Ctrl-P, Up			Previous match from history\n" +
		"Ctrl-N, Down			Next match from history\n" +
		"Ctrl-R				Reverse Search history (Ctrl-S forward, Ctrl-G cancel)\n" +
		"Ctrl-Y				Paste from Yank buffer (Alt-Y to paste next yank instead)\n" +
		"Tab				Next completion\n" +
		"Shift-Tab			(after Tab) Previous completion")
}
