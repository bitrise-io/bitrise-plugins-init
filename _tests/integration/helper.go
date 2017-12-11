package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/bitrise-io/go-utils/command/git"
	"github.com/stretchr/testify/require"
)

func binPath() string {
	return os.Getenv("INTEGRATION_TEST_BINARY_PATH")
}

func gitClone(t *testing.T, dir, uri string) {
	fmt.Printf("cloning into: %s\n", dir)
	g, err := git.New(dir)
	require.NoError(t, err)
	require.NoError(t, g.Clone(uri).Run())
}
