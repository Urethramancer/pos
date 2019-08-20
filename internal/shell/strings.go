package shell

import "github.com/Urethramancer/signor/stringer"

const (
	ErrParseID   = "Error parsing ID"
	ErrConvertID = "Error converting ID"

	ErrAddClient     = "Error adding client"
	ErrUpdateClient  = "Error updating client"
	ErrRemoveClient  = "Error removing client"
	ErrGetClientList = "Error retrieving client list"

	ErrAddContact     = "Error adding contact"
	ErrUpdateContact  = "Error updating contact"
	ErrRemoveContact  = "Error removing contact"
	ErrGetContactList = "Error retrieving contact list"

	ErrAddJob     = "Error adding job"
	ErrEditJob    = "Error editing job"
	ErrRemoveJob  = "Error removing job"
	ErrGetJobList = "Error retrieving job list"
	ErrGetJob     = "Error retrieving job"

	strName        = "Name"
	strEmail       = "E-mail"
	strPhone       = "Phone"
	strCompanyName = "Company name"
	strCompanyID   = "Company ID"
	strAddress     = "Address"
	strVATID       = "VAT ID (org. #)"
)

func (sh *Shell) strPrompt(a, b string) (string, error) {
	buf := stringer.New()
	buf.WriteString(a)
	if b != "" {
		buf.WriteStrings(" [", b, "]")
	}
	buf.WriteString(": ")
	return sh.Prompt(buf.String())
}

func (sh *Shell) intPrompt(s string, n int64) (string, error) {
	buf := stringer.New()
	buf.WriteString(s)
	if n != 0 {
		buf.WriteI(" [", n, "]")
	}
	buf.WriteString(": ")
	return sh.Prompt(buf.String())
}
