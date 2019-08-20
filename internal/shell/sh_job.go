package shell

import (
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/Urethramancer/pos/internal/database"
	"github.com/Urethramancer/signor/stringer"
)

func (sh *Shell) jobCommands(args []string) {
	if len(args) == 0 {
		sh.m("list\t\tList all jobs, optionally only for a client ID.")
		sh.m("add\t\tAdd job. Takes client ID and currency code as arguments.")
		sh.m("edit <id>\tModify a job. Same arguments as adding.")
		sh.m("remove <id...>\tRemove jobs by ID. Will also remove all related tasks.")
		sh.m("show <id>\tShow details and tasks for a job.")
		return
	}

	cmd := args[0]
	args = args[1:]
	switch cmd {
	case "add":
		if len(args) < 2 {
			sh.m("You need to specify a client ID and currency code. Example:\n\tjob add 1 GBP")
		} else {
			sh.addJob(args)
		}

	case "edit":
		if len(args) < 3 {
			sh.m("You need to specify a job ID, client ID and currency code. Example:\n\tjob edit 1 1 GBP")
		} else {
			sh.editJob(args)
		}

	case "remove":
		if len(args) == 0 {
			sh.m("You must specify one or more job IDs to remove.")
			return
		}

		for _, x := range args {
			id, err := strconv.Atoi(x)
			if err != nil {
				sh.e("%s: %s", ErrParseID, err.Error())
				break
			}
			sh.removeJob(int64(id))
		}

	case "list":
		if len(args) < 1 {
			sh.listJobs()
		} else {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				sh.e("%s.", ErrParseID)
				break
			}

			sh.listJobsFor(int64(id))
		}

	case "show":
	}
}

func (sh *Shell) addJob(args []string) {
	client, err := strconv.Atoi(args[0])
	if err != nil {
		sh.e("%s: %s", ErrConvertID, err.Error())
		return
	}

	j := &database.Job{
		Client:   int64(client),
		Currency: args[1],
	}
	id, err := sh.db.AddJob(j)
	if err != nil {
		sh.e("%s: %s", ErrAddJob, err.Error())
		return
	}

	sh.m("Added job with ID %d to client %d.", id, client)
}

func (sh *Shell) editJob(args []string) {
	var err error
	id, err := strconv.Atoi(args[0])
	if err != nil {
		sh.e("%s: %s", ErrConvertID, err.Error())
		return
	}

	client, err := strconv.Atoi(args[1])
	if err != nil {
		sh.e("%s: %s", ErrConvertID, err.Error())
		return
	}

	j := &database.Job{
		ID:       int64(id),
		Client:   int64(client),
		Currency: args[1],
	}
	err = sh.db.UpdateJob(j)
	if err != nil {
		sh.e("%s: %s", ErrEditJob, err.Error())
		return
	}

	sh.m("Updated job %d for client %d.", id, client)
}

func (sh *Shell) removeJob(id int64) {
	err := sh.db.RemoveJob(id)
	if err != nil {
		sh.e("%s: %s", ErrRemoveJob, err.Error())
		return
	}

	sh.m("Removed job %d.", id)
}

func (sh *Shell) listJobs() {
	list, err := sh.db.GetAllJobs()
	if err != nil {
		sh.e("%s: %s", ErrGetJobList, err.Error())
		return
	}

	sh.printJobs(list)
}

func (sh *Shell) listJobsFor(id int64) {
	list, err := sh.db.GetJobsFor(id)
	if err != nil {
		sh.e("%s: %s", ErrGetJobList, err.Error())
		return
	}

	sh.printJobs(list)
}

func (sh *Shell) showJob(id int64) {
	j, err := sh.db.GetJob(id)
	if err != nil {
		sh.e("%s: %s", ErrGetJob, err.Error())
		return
	}

	sh.printJob(j)
}

func (sh *Shell) printJob(j *database.Job) {
	sh.m("Job %d for %s (%d): %d tasks")
}

func (sh *Shell) printJobs(list []*database.Job) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	s := stringer.New()
	s.WriteString("ID\tClient\tCurrency\tCost\tCreated\n")
	for _, j := range list {
		company := "Unknown"
		client := sh.db.GetClient(j.Client)
		if client != nil {
			company = client.Company
		}
		_, err := s.WriteI(
			j.ID,
			"\t", j.Client, " (", company, ")",
			"\t", j.Currency,
			"\t", j.Cost,
			"\t", j.Created.String(),
			"\n",
		)

		if err != nil {
			sh.e("Error printing to stdout: %s", err.Error())
			return
		}
	}
	_, _ = tw.Write([]byte(s.String()))
	tw.Flush()
}
