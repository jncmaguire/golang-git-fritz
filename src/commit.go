package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	// delimeter separates the Conventional Commits header, body, and footer.
	delimeter string = "\n\n"
	// commentExpr is the expression to match comments in a git commit.
	commentExpr string = `#.*\n?`
	// breakingChange is the string that must prepend a breaking change summary.
	breakingChange string = "BREAKING CHANGE: "
	// breakingSymbol is the symbol that can be added to a commit type in order to denote that it's a breaking change.
	breakingSymbol string = "!"
)

var (
	// these types are required, regardless of how the user wants to personalize things.
	requiredKinds = []string{"fix", "feat"}
)

// commitValidator holds the regexp expressions to match various components of a commit message.
type commitValidator struct {
	delimeter,
	header,
	breaking,
	breakingInBody,
	ticket,
	comments *regexp.Regexp
}

// commit represents a git commit.
type commit struct {
	file,
	kind,
	scope,
	subject,
	ticket,
	header,
	body,
	footer,
	message string
	hasTicket,
	isBreaking bool
	validator commitValidator
}

// newCommit creates a new commit. Very little validation happens at this stage.
func newCommit(location string, kind string, scope string, subject string, ticket string, isBreaking bool) (comm *commit) {
	comm.file = location
	comm.scope = scope
	comm.subject = subject
	comm.ticket = ticket

	comm.kind = strings.ToLower(kind)

	comm.isBreaking = isBreaking || strings.Contains(kind, breakingSymbol)
	comm.hasTicket = len(ticket) > 0

	return comm
}

// write builds a commit message, sets up the validators based on the user configuration, validates the commit, and saves it (if valid).
func (comm *commit) write(cfg config) (err error) {
	// build a commit message
	err = comm.build()

	if err == nil {
		err = comm.setupValidators(cfg.otherKinds, cfg.scopes, cfg.ticketFormat)
	}

	if err == nil {
		err = comm.validate()
	}

	if err == nil {
		err = comm.save()
	}
	return err
}

// parse reads the commit message into the struct, sets up the validators based on the user configuration, and validates the commit message.
func (comm *commit) parse(cfg config) (err error) {
	err = comm.read()

	if err == nil {
		err = comm.setupValidators(cfg.otherKinds, cfg.scopes, cfg.ticketFormat)
	}

	if err == nil {
		err = comm.validate()
	}

	return err
}

// setupValidators creates the regexp matchers to be used in validation. It can take user configured values.
func (comm *commit) setupValidators(otherKinds []string, userScopes []string, ticketFormat string) (err error) {

	// compile the regexp pattern for types
	kinds := []string{}
	for _, kind := range append(requiredKinds, otherKinds...) {
		kinds = append(kinds, fmt.Sprintf("(%s)", kind))
	}

	kindExpr := strings.Join(kinds, "|") + breakingSymbol + "?"

	// compile the regexp pattern for scopes; default is any word of any length
	scopeExpr := `\w+`
	if len(userScopes) > 0 {
		scopes := []string{}
		for _, scope := range userScopes {
			scopes = append(scopes, fmt.Sprintf("(%s)", scope))
		}

		scopeExpr = strings.Join(scopes, "|")
	}

	// compile the delimiter into a regexp pattern
	comm.validator.delimeter, err = regexp.Compile(delimeter)

	if err == nil {
		// compile the regexp pattern for a header, based on the above
		comm.validator.header, err = regexp.Compile(fmt.Sprintf(`^(%s)*%s(\(%s\))?: .+`, commentExpr, kindExpr, scopeExpr))
	}

	// if it's a breaking commit, compile the regexp pattern to find 'breaking commit' notes where 'breaking commit' notes are allowed to appear.
	if err == nil && comm.isBreaking {
		breaking := breakingChange + `.+`
		comm.validator.breaking, err = regexp.Compile(`(^|\n)` + breaking)

		if err == nil {
			comm.validator.breakingInBody, err = regexp.Compile(`^` + breaking)
		}
	}

	// if a ticket was passed with the ticket, compile the regexp pattern to find the notes mentioning the ticket.
	if err == nil && comm.hasTicket {
		comm.validator.ticket, err = regexp.Compile(fmt.Sprintf(`(^|\n).+\b%s\b`, ticketFormat))
	}

	return err
}

// read reads the commit message from the commit file.
func (comm *commit) read() error {
	message, err := ioutil.ReadFile(comm.file)

	if len(message) == 0 {
		err = fmt.Errorf("commit message is empty")
	}

	if err == nil {
		// split the commit message into sections using the delimeter noted in Conventional Commits
		sections := comm.validator.delimeter.Split(string(message), -1)

		// the header is always the first in the index
		comm.header = sections[0]

		total := len(sections)

		beg := 1
		end := total

		if total == 1 && (comm.hasTicket || comm.isBreaking) {
			err = fmt.Errorf("body or footer is required but missing")
		}

		// if there's more than 1 section, then there's a body and / or footer
		if err == nil && total > 1 {
			end--

			// the footer is always the last in the index
			comm.footer = sections[end]

			if comm.hasTicket {
				// the body is everything in between
				comm.body = strings.Join(sections[beg:end], delimeter)
			}
		}
	}

	return err
}

// build creates the commit message based on the struct values.
func (comm *commit) build() (err error) {
	scope := ""
	if len(comm.scope) > 0 {
		scope = fmt.Sprintf("(%s)", comm.scope)
	}

	comm.header = fmt.Sprintf("%s%s: %s", comm.kind, scope, comm.subject)

	if comm.isBreaking {
		comm.body = delimeter + breakingChange
	}

	if comm.hasTicket {
		comm.footer = delimeter + "closes " + comm.ticket
	}

	comm.message = comm.header + comm.body + comm.footer

	return err
}

// validate checks that the commit message is valid based on Conventional Commits and the user configuration.
func (comm *commit) validate() (err error) {
	if !comm.validator.header.MatchString(comm.header) {
		err = fmt.Errorf("header is invalid")
	} else if len([]rune(comm.header)) > 50 {
		err = fmt.Errorf("header cannot be longer than 50 characters")
	}

	if comm.isBreaking && !(comm.validator.breaking.MatchString(comm.footer) || comm.validator.breakingInBody.MatchString(comm.body)) {
		err = fmt.Errorf("%s\nbreaking message missing", err)
	}

	if comm.hasTicket && !comm.validator.ticket.MatchString(comm.footer) {
		err = fmt.Errorf("%s\nticket missing", err)
	}

	// if there's an error, print it in a comment at the top
	if err != nil {
		comm.header = `# ERROR: %v\n` + comm.header
	}

	return err
}

// save writes the prepared commit message to the commit file.
func (comm *commit) save() error {
	return ioutil.WriteFile(comm.file, []byte(comm.message), 0644)
}
