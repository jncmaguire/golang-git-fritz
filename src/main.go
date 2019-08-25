package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	// subcommands
	operSetup   string = "setup"
	operCleanup string = "cleanup"

	// feature subcommands, subcommands that are for using the program's features
	operCommitPrep     string = "commit-prep"
	operCommitValidate string = "commit-validate"
)

var (
	arg  args
	sets map[string]*flag.FlagSet
)

func init() {
	sets = arg.setFlags()
}

func main() {

	// any error should cause a panic; we can't continue past an error
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	// check that we're called in the same directory as a .git folder
	err := checkForGit()

	// parse flags and act based on a subcommand

	if err == nil {
		err = arg.parseFlags()
	}

	if err == nil {

		// fritz.toml is the required name for the config file. It must be in the same directory as where the .git folder is.
		configFile := "fritz.toml"

		switch arg.subcommand {
		case operSetup:
			err = setup(configFile, arg.features)
		case operCleanup:
			err = cleanup(configFile)
		case operCommitPrep:
			var cfg config
			comm := newCommit(arg.location, arg.kind, arg.scope, arg.subject, arg.ticket, arg.isBreaking)

			// load config file to personalize validation
			cfg, err = loadConfig(configFile)
			err = comm.write(cfg)
		case operCommitValidate:
			var cfg config
			comm := commit{file: arg.location}

			// load config file to personalize validation
			cfg, err = loadConfig(configFile)
			err = comm.parse(cfg)
		default:
			err = fmt.Errorf("Not a valid command")
		}
	}

	if err != nil {
		panic(err.Error())
	}
}
