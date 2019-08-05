package main

import (
	"github.com/Urethramancer/signor/files"
)

type Translation struct {
	table map[string]map[string]string
	lang  string
}

// Get translation.
func (tr *Translation) Get(s string) string {
	x, ok := tr.table[s][tr.lang]
	if !ok {
		return s
	}

	return x
}

func loadTranslation(lang string) (*Translation, error) {
	t := Translation{
		table: make(map[string]map[string]string),
		lang:  lang,
	}
	err := files.LoadJSON("strings.json", &t.table)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
