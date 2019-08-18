package shell

import (
	"io"
	"os"
	"sort"
	"strings"

	"github.com/Urethramancer/cross"
	"github.com/Urethramancer/pos/internal/config"
	"github.com/Urethramancer/pos/internal/database"
	"github.com/Urethramancer/signor/log"
	"github.com/peterh/liner"
)

// Shell for interactive management.
type Shell struct {
	cfg *config.Config
	db  *database.Invoices
	*liner.State

	history string

	m func(string, ...interface{})
	e func(string, ...interface{})
}

// New creates an interactive shell structure and loads the configuration.
func New() (*Shell, error) {
	sh := &Shell{}
	sh.cfg = &config.Config{
		Name:          "",
		Email:         "",
		Host:          "localhost",
		Port:          "5432",
		DBName:        "invoices",
		Username:      "postgres",
		Password:      "postgres",
		Company:       "",
		CompanyID:     "",
		VAT:           "0",
		InvoicePrefix: "",
	}

	cross.SetConfigPath("pos")
	fn := cross.ConfigName("config.json")
	err := sh.cfg.Load(fn)
	if err != nil {
		log.Default.Msg("Creating new configuration file.")
		err = sh.cfg.Save(fn)
		if err != nil {
			return nil, err
		}
	}

	sh.db, err = database.Open(sh.cfg.Host, sh.cfg.Port, sh.cfg.Username, sh.cfg.Password, sh.cfg.DBName, false)
	if err != nil {
		return nil, err
	}

	sh.history = cross.ConfigName("history")
	sh.m = log.Default.Msg
	sh.e = log.Default.Err
	return sh, nil
}

// Close the database and do any other cleanup required.
func (sh *Shell) Close() {
	if sh.db != nil {
		sh.db.Close()
	}

	if sh.State != nil {
		sh.State.Close()
	}
}

// LoadHistory from file.
func (sh *Shell) LoadHistory() {
	sh.State = liner.NewLiner()
	// sh.SetHistoryLimit(100)
	f, err := os.Open(sh.history)
	if err != nil {
		return
	} else {
		defer f.Close()
		n, err := sh.ReadHistory(f)
		if err != nil {
			sh.e("Couldn't load history: %s", err.Error())
			return
		} else {
			sh.m("Loaded %d history items.", n)
		}
	}
}

// SaveHistory to file.
func (sh *Shell) SaveHistory() {
	f, err := os.OpenFile(sh.history, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		sh.e("Couldn't open history file %s: %s", sh.history, err.Error())
	}

	defer f.Close()
	_, err = sh.WriteHistory(f)
	if err != nil {
		sh.e("Couldn't write history: %s", err.Error())
	}
}

// Run the shell until exited or some fatal error happens.
func (sh *Shell) Run() error {
	sh.LoadHistory()
	defer sh.SaveHistory()
	sh.m("Ctrl-D to quit, Tab for command completion, ? for a list of commands, ?? for a full shortcut list.")

	commands := []string{
		"setup",
		"client", "contact",
		"job", "task", "invoice",
	}
	sort.Strings(commands)
	sh.SetCompleter(func(line string) (c []string) {
		for _, n := range commands {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	running := true
	for running {
		x, err := sh.Prompt("> ")
		if err != nil {
			if err == io.EOF {
				sh.m("")
				return nil
			}

			return err
		}

		if x != "" {
			sh.AppendHistory(x)
		}

		cmd, args := splitCommand(x)
		switch cmd {
		case "?":
			sh.m("%s", Tablify(commands, 5))
		case "??":
			showShortcuts()
		case "setup":
			err = sh.Setup()
			if err != nil {
				if err == io.EOF {
					sh.m("Setup aborted.")
				} else {
					sh.e("Error during setup: %s", err.Error())
				}
			}

		case "client":
			sh.clientCommands(args)

		case "contact":
			sh.contactCommands(args)

		case "task":

		case "job":
			sh.jobCommands(args)

		case "invoice":
		default:
			if cmd != "" {
				sh.m("Unknown command '%s'.", cmd)
			}
		}
	}

	return nil
}

func splitCommand(s string) (string, []string) {
	a := strings.Split(s, " ")
	if len(a) == 1 {
		return a[0], nil
	}

	return a[0], a[1:]
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
