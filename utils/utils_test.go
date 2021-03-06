package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TouchFile(t *testing.T, fname string) *os.File {
	f, err := os.Create(fname)
	assert.Nil(t, err)
	return f
}

func Test_FileExists_ReturnsTrue_IfFileExists(t *testing.T) {
	assert.True(t, FileExists(os.Getenv("HOME")))
}

func Test_FileExists_ReturnsFalse_IfFileDoesNotExist(t *testing.T) {
	assert.False(t, FileExists("/foo"))
}

func Test_Canonicalize(t *testing.T) {
	home := os.Getenv("HOME")
	assert.NotEmpty(t, home)
	var testCases = []struct {
		path         string
		expandedPath string
	}{
		{"~", home},
		{"~/foo/bar", fmt.Sprintf("%s/foo/bar", home)},
		{"/foo/bar", "/foo/bar"},
		{"", ""},
	}
	for _, testCase := range testCases {
		assert.Equal(t, testCase.expandedPath, Canonicalize(testCase.path))
	}
}

func Test_IsDirectory_ReturnsTrue_WhenPathIsDirectory(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	result, err := IsDirectory(path)
	assert.True(t, result)
	assert.Nil(t, err)
}

func Test_IsDirectory_ReturnsFalse_WhenPathIsFile(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	fname := filepath.Join(path, "test_file")
	TouchFile(t, fname).Close()

	result, err := IsDirectory(fname)
	assert.False(t, result)
	assert.Nil(t, err)
}

func Test_ConcatFiles(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	f_a := filepath.Join(path, "a")
	f := TouchFile(t, f_a)
	f.Write([]byte("aaa"))
	f.Close()

	f_b := filepath.Join(path, "b")
	f = TouchFile(t, f_b)
	f.Write([]byte("bbb"))
	f.Close()

	f_c := filepath.Join(path, "c")

	err = ConcatFiles(f_c, []string{f_a, f_b}, path)
	assert.Nil(t, err)
	contents, err := ioutil.ReadFile(f_c)
	assert.Equal(t, "aaabbb", string(contents))
}
