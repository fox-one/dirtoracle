package cmd

func getArg(args []string, idx int) (string, bool) {
	if idx < len(args) {
		return args[idx], true
	}

	return "", false
}
