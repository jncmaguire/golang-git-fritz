package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path"
)

const (
	gitPath   string = ".git"  // name of the repository-level git folder
	hooksPath string = "hooks" // name of the hooks folder in the gitPath
)

// config holds the user-personalized values used in commit validation.
type config struct {
	hooks        []string
	otherKinds   []string `toml:"otherTypes"`
	scopes       []string
	ticketFormat string
}

// checkForGit makes sure that the program is being called from the same directory as where the .git folder is.
func checkForGit() (err error) {
	var workingDir string
	workingDir, err = os.Getwd()
	if err == nil {
		gitPath := path.Join(workingDir, gitPath, hooksPath)
		_, gitPathErr := os.Stat(gitPath)
		if os.IsNotExist(gitPathErr) {
			err = fmt.Errorf("call program from same level as .git directory in repo")
		}
	}
	return err
}

// loadConfig loads the user configuration values into a struct.
func loadConfig(configFileName string) (cfg config, err error) {
	var file []byte
	file, err = ioutil.ReadFile(configFileName)

	if err == nil {
		err = toml.Unmarshal(file, cfg)
	}

	return cfg, err
}

// setup sets a git repository up to be used with the program. It creates the configuration file as well as the githooks.
func setup(configFileName string, features []string) (err error) {
	var (
		file          []byte
		defaultConfig = config{
			otherKinds: []string{"improvement"},
		}
	)

	if len(features) == 0 {
		err = fmt.Errorf("need to select at least 1 feature for setup")
	}

	if err == nil {
		// add the appropriate hooks to the configuration file struct
		defaultConfig.hooks = getHooks(features)
		file, err = toml.Marshal(defaultConfig)
	}

	if err == nil {
		// create the scripts
		for _, hook := range defaultConfig.hooks {
			script := scriptByHook[hook]
			err = ioutil.WriteFile(path.Join(gitPath, hooksPath, hook), []byte(script.build()), 0644)
		}
	}

	if err == nil {
		// create the configuration file
		err = ioutil.WriteFile(configFileName, file, 0644)
	}

	return err
}

// cleanup removes the program from the repository. It removes any created hooks and deletes the configuration file.
func cleanup(configFileName string) error {
	// open config
	cfg, err := loadConfig(configFileName)

	if err == nil {
		// remove created hooks
		for _, hook := range cfg.hooks {
			err = os.Remove(path.Join(gitPath, hooksPath, hook))

			if err != nil {
				break
			}
		}
	}

	if err == nil {
		// delete config file
		err = os.Remove(configFileName)
	}

	return err
}
