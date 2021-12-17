package main

import (
	"fmt"

	"github.com/fox-one/dirtoracle/apps/mvm/cmd"
)

var (
	version string
	commit  string
)

func main() {
	ver := fmt.Sprintf("%s (%s)", version, commit)
	cmd.Execute(ver)
}
