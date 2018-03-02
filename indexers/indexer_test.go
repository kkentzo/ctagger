package indexers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kkentzo/tagger/utils"
	"github.com/kkentzo/tagger/watchers"
	"github.com/stretchr/testify/assert"
)

func TouchFile(t *testing.T, fname string) *os.File {
	f, err := os.Create(fname)
	assert.Nil(t, err)
	return f
}

func CheckGenericArguments(t *testing.T, args []string) {
	assert.Contains(t, args, "-R")
	assert.Contains(t, args, "-e")
	assert.Contains(t, args, "--exclude=.git")
}

func Test_Indexer_Deserialization(t *testing.T) {
	t.Skip("TODO")
}

func Test_Indexer_DefaultIndexer(t *testing.T) {
	indexer := DefaultIndexer()
	assert.Equal(t, "ctags", indexer.Program)
	assert.Contains(t, indexer.Args, "-R")
	assert.Contains(t, indexer.Args, "-e")
	assert.Equal(t, "TAGS", indexer.TagFilePrefix)
	assert.Equal(t, Generic, indexer.Type)
	assert.Contains(t, indexer.ExcludeDirs, ".git")
}

func Test_Indexer_Index_ShouldTriggerCommand(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	indexer := &Indexer{
		Program: "touch",
		Args:    []string{"aaa"},
		Type:    Generic,
	}

	indexer.Index(path, watchers.Event{})
	assert.True(t, utils.FileExists(filepath.Join(path, "aaa")))
}

func Test_Indexer_CreateWatcher_ShouldReturnAWatcher(t *testing.T) {
	indexer := &Indexer{
		MaxPeriod: 2 * time.Second,
		Type:      Rvm,
	}
	watcher := indexer.CreateWatcher("foo").(*watchers.Watcher)
	defer watcher.Close()

	assert.Equal(t, "foo", watcher.Root)
	assert.Equal(t, "Gemfile.lock", watcher.SpecialFile)
	assert.Equal(t, 2*time.Second, watcher.MaxPeriod)
}

func Test_Indexer_GetGenericArguments(t *testing.T) {
	indexer := DefaultIndexer()
	args := indexer.GetGenericArguments("foo")
	CheckGenericArguments(t, args)
	assert.Equal(t, 3, len(args))
}

func Test_Indexer_GetProjectArguments(t *testing.T) {
	indexer := DefaultIndexer()
	args := indexer.GetProjectArguments("foo")
	CheckGenericArguments(t, args)
	assert.Contains(t, args, "-f TAGS.project")
	assert.Equal(t, ".", args[len(args)-1])
}

func Test_Indexer_GetGemsetArguments_WhenIndexerIsRvm(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	// Prepare rvm-specific files
	TouchFile(t, filepath.Join(path, "Gemfile")).Close()
	f := TouchFile(t, filepath.Join(path, ".ruby-version"))
	f.Write([]byte("2.1.3"))
	f.Close()
	f = TouchFile(t, filepath.Join(path, ".ruby-gemset"))
	f.Write([]byte("foo"))
	f.Close()

	indexer := DefaultIndexer()
	indexer.Type = Rvm
	args := indexer.GetGemsetArguments(path)
	CheckGenericArguments(t, args)
	gp, err := rvmGemsetPath(path)
	assert.Nil(t, err)
	assert.Contains(t, args, "-f TAGS.gemset")
	assert.Equal(t, gp, args[len(args)-1])
}

func Test_Indexer_GetGemsetArguments_WhenIndexerIsNotRvm(t *testing.T) {
	indexer := DefaultIndexer()
	indexer.Type = Rvm
	args := indexer.GetGemsetArguments("foo")
	assert.Empty(t, args)
}

func Test_Indexer_GetTagFileNameForGemset(t *testing.T) {
	indexer := DefaultIndexer()
	assert.Equal(t, "aaa/TAGS.gemset", indexer.GetTagFileNameForGemset("aaa"))

}

func Test_Indexer_GemsetTagFileExists_ReturnsTrue_WhenTagFileExists(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	TouchFile(t, filepath.Join(path, "TAGS.gemset")).Close()

	indexer := DefaultIndexer()
	assert.True(t, indexer.GemsetTagFileExists(path))
}

func Test_Indexer_GemsetTagFileExists_ReturnsFalse_WhenTagFileDoesNotExist(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	indexer := DefaultIndexer()
	assert.False(t, indexer.GemsetTagFileExists(path))
}

func Test_Indexer_GetTagFileNameForProject(t *testing.T) {
	indexer := DefaultIndexer()
	assert.Equal(t, "aaa/TAGS.project", indexer.GetTagFileNameForProject("aaa"))
}

func Test_Indexer_ProjectTagFileExists_ReturnsTrue_WhenTagFileExists(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	TouchFile(t, filepath.Join(path, "TAGS.project")).Close()

	indexer := DefaultIndexer()
	assert.True(t, indexer.ProjectTagFileExists(path))
}

func Test_Indexer_ProjectTagFileExists_ReturnsFalse_WhenTagFileDoesNotExist(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	indexer := DefaultIndexer()
	assert.False(t, indexer.ProjectTagFileExists(path))
}