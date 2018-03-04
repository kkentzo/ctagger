package indexers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkentzo/tagger/utils"
)

func isRuby(root string) bool {
	return utils.FileExists(filepath.Join(root, "Gemfile"))
}

func isRvm(root string) bool {
	return isRuby(root) &&
		utils.FileExists(filepath.Join(root, ".ruby-version")) &&
		utils.FileExists(filepath.Join(root, ".ruby-gemset"))
}

func rubyVersion(root string) (string, error) {
	rv, err := ioutil.ReadFile(filepath.Join(root, ".ruby-version"))
	if err != nil {
		return "", err
	} else {
		return strings.TrimSpace(string(rv)), nil
	}
}

func rubyGemset(root string) (string, error) {
	rg, err := ioutil.ReadFile(filepath.Join(root, ".ruby-gemset"))
	if err != nil {
		return "", err
	} else {
		return strings.TrimSpace(string(rg)), nil
	}
}

func rvmGemsetPathFromFiles(root string) (string, error) {
	rv, err := rubyVersion(root)
	if err != nil {
		return "", err
	}
	rg, err := rubyGemset(root)
	if err != nil {
		return "", err
	}
	path := filepath.Join(os.Getenv("HOME"),
		fmt.Sprintf(".rvm/gems/%s@%s/gems", rv, rg))
	return path, nil

}

// TODO: Switch to this implementation
func rvmGemsetPathFromRvm(root string) (string, error) {
	cmd := "/bin/bash"
	args := []string{
		"-c",
		"source \"$HOME/.rvm/scripts/rvm\"; cd .; rvm gemset gemdir"}
	out, err := utils.ExecInPath(cmd, args, root)
	if err != nil {
		gemset := strings.TrimSpace(string(out))
		return filepath.Join(gemset, "gems"), nil
	} else {
		return "", errors.New(fmt.Sprint(string(out), err.Error()))
	}
}
