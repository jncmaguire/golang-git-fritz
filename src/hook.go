package main

const (
	// official githook names
	hookPrepareCommitMsg string = "prepare-commit-msg"
	hookCommitMsg        string = "commit-msg"
)

// maps a feature subcommand to the hook during which it occurs
var hooksByFeature = map[string]string{
	operCommitPrep:     hookPrepareCommitMsg,
	operCommitValidate: hookCommitMsg,
}

// getHooks takes a list of features and outputs the hooks we must create for them.
func getHooks(features []string) (hooks []string) {
	for _, feature := range features {
		hooks = append(hooks, hooksByFeature[feature])
	}
	return hooks
}
