package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-plugins-init/version"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

var helpStr = fmt.Sprintf(`NAME:
   bitrise-plugins-init - Init bitrise config

USAGE:
   bitrise-plugins-init [global options] command [command options] [arguments...]

VERSION:
   %s

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --minimal      create empty bitrise config and secrets
   --help, -h     show help
   --version, -v  print the version`, version.VERSION)

func Test_HelpTest(t *testing.T) {
	t.Log("help command")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPath(), "help")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
		require.Equal(t, helpStr, out)
	}

	t.Log("help short command")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPath(), "h")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
		require.Equal(t, helpStr, out)
	}

	t.Log("help flag")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPath(), "--help")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
		require.Equal(t, helpStr, out)
	}

	t.Log("help short flag")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPath(), "-h")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
		require.Equal(t, helpStr, out)
	}
}
