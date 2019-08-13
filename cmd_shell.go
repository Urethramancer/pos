package main

import (
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/peterh/liner"

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

	m := log.Default.Msg
	m("Ctrl-D to quit, Tab for command completion, ? for a list of commands, ?? for a full shortcut list.")
	line := liner.NewLiner()

	commands := []string{"setup",
		"clients", "cl",
		"invoices", "inv",
		"jobs", "tasks",
	}
	sort.Strings(commands)
	line.SetCompleter(func(line string) (c []string) {
		for _, n := range commands {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	defer line.Close()
	running := true
	for running {
		x, err := line.Prompt("> ")
		if err != nil {
			if err == io.EOF {
				m("")
				return nil
			}

			return err
		}

		line.AppendHistory(x)
		switch x {
		case "?":
			m("%s", Tablify(commands, 5))
		case "??":
			showShortcuts()
		}
	}

	return nil
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
