package main

import (
	"flag"
	"os"
)

// args represents all the possible command line arguments the program has.
type args struct {
	features []string
	subcommand,
	location,
	kind,
	scope,
	subject,
	ticket string
	isBreaking bool
}

// setFlags sets the flags for each subcommand, such that the values are loaded into the args object. It returns a map of FlagSets keyed by subcommand name.
func (a *args) setFlags() (setBySubcommand map[string]*flag.
	FlagSet) {
	setBySubcommand = make(map[string]*flag.FlagSet)
	subcommands := []string{operSetup, operCleanup, operCommitPrep, operCommitValidate}
	for _, subcommand := range subcommands {
		set := flag.NewFlagSet(subcommand, flag.ExitOnError)

		if subcommand == operCommitPrep {
			set.StringVar(&a.kind, "type", "", "kind of commit (e.g. 'fix')")
			set.StringVar(&a.scope, "scope", "", "feature it relates to (e.g. 'header-component')")
			set.StringVar(&a.subject, "subject", "", "short description of what the commit does")
			set.StringVar(&a.ticket, "ticket", "", "ticket for commit")
			set.BoolVar(&a.isBreaking, "breaking", false, "symbol for if the commit has breaking changes")
		} else {
			set.StringVar(&a.location, "location", "", "file location of the commit")
		}
		setBySubcommand[subcommand] = set
	}
	return setBySubcommand
}

// parseFlags uses the subcommand name to parse the remainder of the command line arguments appropriately.
func (a *args) parseFlags() (err error) {
	subcommand := os.Args[1]

	if set, ok := sets[subcommand]; ok {
		err = set.Parse(os.Args[2:])
	} else if subcommand == `setup` {
		arg.features = os.Args[2:]
	}

	if err == nil {
		arg.subcommand = subcommand
	}

	return err
}
