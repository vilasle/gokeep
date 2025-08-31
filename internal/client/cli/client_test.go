package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vilasle/gokeep/internal/client"
)

func TestCommandLineClient_New(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	testdata := filepath.Join(pwd, "testdata")
	os.RemoveAll(testdata)

	err = os.MkdirAll(testdata, 0755)
	require.NoError(t, err)

	err = client.CreateNewConfiguration(testdata, ":55654", "local.db")
	require.NoError(t, err)

	config, err := client.GetCurrentConfiguration(testdata)
	require.NoError(t, err)

	c, err := NewClient(config)
	require.NoError(t, c.Close())
}
