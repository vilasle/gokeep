package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateConfiguration(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	testdata := filepath.Join(pwd, "testdata")

	grpcAddr := ":0"
	db := filepath.Join(testdata, "local.db")

	os.RemoveAll(testdata)

	err = os.Mkdir(testdata, 0755)
	require.NoError(t, err)

	err = CreateNewConfiguration(testdata, grpcAddr, db)
	require.NoError(t, err)

	ls, err := os.ReadDir(testdata)
	require.NoError(t, err)

	assert.NotEmpty(t, ls)
}

func Test_GetConfiguration(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	testdata := filepath.Join(pwd, "testdata")

	config, err := GetCurrentConfiguration(testdata)
	require.NoError(t, err)

	require.NotEmpty(t, config)

	assert.Equal(t, filepath.Join(testdata, "cert"), config.Certificate.Path)
	assert.Equal(t, filepath.Join(testdata, "config.yaml"), config.Config.Path)
	assert.Equal(t, testdata, config.ConfigDirectory.Path)
	assert.Equal(t, testdata, config.Credentials.Path)
	assert.Equal(t, filepath.Join(testdata, "upload"), config.UploadDirectory.Path)

}
