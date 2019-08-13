package main

import (
	"strings"

	"github.com/Urethramancer/signor/stringer"
)

// Tablify in to w words per line and output tab-delimited.
func Tablify(in []string, w int) string {
	c := 0
	max := 0
	for _, x := range in {
		if len(x) > max {
			max = len(x)
		}
	}

	tabs := max/8 + 1
	buf := stringer.New()
	var t string
	for _, x := range in {
		t = strings.Repeat("\t", tabs-(len(x)/8))
		if c == 0 {
			buf.WriteString("\t")
		}
		if c < w {
			buf.WriteStrings(x, t)
		} else {
			buf.WriteStrings("\n\t", x, t)
			c = 0
		}
		c++
	}
	return buf.String()
}
