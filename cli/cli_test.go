package cli

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

const ignorePattern = ".bitrise.secrets.yml"

func Test_GitignoreTest(t *testing.T) {
	t.Log("create .gitignore with pattern when .gitignore does not exist")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		gitignorePath := path.Join(tmpDir, ".gitignore")
		err = gitignore(ignorePattern, gitignorePath)
		require.NoError(t, err)

		contents, err := ioutil.ReadFile(gitignorePath)
		require.NoError(t, err)

		require.Equal(t, ".bitrise.secrets.yml", string(contents))

	}

	t.Log("write on last line in .gitignore when file ends with new line")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		gitignorePath := path.Join(tmpDir, ".gitignore")
		content := []byte("node_modules\n")
		err = ioutil.WriteFile(gitignorePath, content, 0644)
		require.NoError(t, err)

		err = gitignore(ignorePattern, gitignorePath)
		require.NoError(t, err)

		contents, err := ioutil.ReadFile(gitignorePath)
		require.NoError(t, err)
		require.Equal(t, "node_modules\n.bitrise.secrets.yml", string(contents))
	}

	t.Log("append to new line in .gitignore when file does not end with new line")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		gitignorePath := path.Join(tmpDir, ".gitignore")
		content := []byte("node_modules")
		err = ioutil.WriteFile(gitignorePath, content, 0644)
		require.NoError(t, err)

		err = gitignore(ignorePattern, gitignorePath)
		require.NoError(t, err)

		contents, err := ioutil.ReadFile(gitignorePath)
		require.NoError(t, err)
		require.Equal(t, "node_modules\n.bitrise.secrets.yml", string(contents))
	}

	t.Log("do not add .bitrise.secrets.yml to .gitignore if it's already there")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		gitignorePath := path.Join(tmpDir, ".gitignore")
		content := []byte(".bitrise.secrets.yml")
		err = ioutil.WriteFile(gitignorePath, content, 0644)
		require.NoError(t, err)

		err = gitignore(ignorePattern, gitignorePath)
		require.NoError(t, err)

		contents, err := ioutil.ReadFile(gitignorePath)
		require.NoError(t, err)
		require.Equal(t, ".bitrise.secrets.yml", string(contents))
	}
}
