package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
)

// CmdInvoiceList options.
type CmdInvoiceList struct {
	opt.DefaultHelp
}

// Run list
func (cmd *CmdInvoiceList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	m := log.Default.Msg
	tr, err := loadTranslation("no")
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile("tpl/index.html")
	if err != nil {
		return err
	}

	tpl, err := template.New("invoice").Parse(string(file))
	if err != nil {
		return err
	}

	data := InvoiceTemplate{
		InvoiceTitle: tr.Get("INVOICE"),
	}
	err = tpl.Execute(os.Stdout, data)
	if err != nil {
		return err
	}

	m("Done.")
	return nil
}
