package integration

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func Test_InitTest(t *testing.T) {
	t.Log("init --minimal - platform independent - SHOULD SUCCEED")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPath(), "--minimal")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
	}

	t.Log(".bitrise.secrets.yml added to .gitignore when command succeeds")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		exists, err := pathutil.IsPathExists(tmpDir + "/.gitignore")
		require.NoError(t, err)
		require.False(t, exists, fmt.Sprintf(".gitignore file should not exist at %s", tmpDir+"/.gitignore"))

		cmd := command.New(binPath(), "--minimal")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		content, err := ioutil.ReadFile(tmpDir + "/.gitignore")
		require.NoError(t, err, out)

		require.True(t, strings.Contains(string(content), ".bitrise.secrets.yml"))

	}

	t.Log("init --minimal - bitrise.yml already exists - SHOULD FAIL")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		_, err = os.Create(filepath.Join(tmpDir, "bitrise.yml"))
		require.NoError(t, err)

		cmd := command.New(binPath())
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.EqualError(t, err, "exit status 1", out)
	}

	t.Log("init --minimal - .bitrise.secrets.yml already exists - SHOULD FAIL")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		_, err = os.Create(filepath.Join(tmpDir, ".bitrise.secrets.yml"))
		require.NoError(t, err)

		cmd := command.New(binPath())
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.EqualError(t, err, "exit status 1", out)
	}

	t.Log("init - cordova platform detected - SHOULD SUCCEED")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-cordova-with-jasmine.git"
		gitClone(t, tmpDir, sampleAppURL)

		cmd := command.New(binPath())
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.EqualError(t, err, "exit status 1", out)
	}

	t.Log("init - no platform detected - SHOULD FAIL")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPath())
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.EqualError(t, err, "exit status 1", out)
	}

	t.Log("init - bitrise.yml already exists - SHOULD FAIL")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		_, err = os.Create(filepath.Join(tmpDir, "bitrise.yml"))
		require.NoError(t, err)

		cmd := command.New(binPath())
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.EqualError(t, err, "exit status 1", out)
	}

	t.Log("init - .bitrise.secrets.yml already exists - SHOULD FAIL")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		_, err = os.Create(filepath.Join(tmpDir, ".bitrise.secrets.yml"))
		require.NoError(t, err)

		cmd := command.New(binPath())
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.EqualError(t, err, "exit status 1", out)
	}
}
