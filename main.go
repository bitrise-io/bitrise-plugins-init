package main

import (
	"github.com/bitrise-io/bitrise-plugins-init/cli"
	_ "github.com/bitrise-io/go-utils/command/git"
	_ "github.com/stretchr/testify/require"
)

func main() {
	cli.Run()
}
