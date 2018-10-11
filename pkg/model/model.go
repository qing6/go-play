package model

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type PackageReplacer struct {
	ImportPath    string
	GitRemoteRepo string
}

func (pkg PackageReplacer) ExistOnDisk() (dir string, exist bool, isGitRepo bool) {
	dir = fmt.Sprintf("%s/src/%s", gopath, pkg.ImportPath)
	if _, err := os.Stat(dir); err != nil {
		return
	}
	exist = true
	dirDotGit := strings.Join([]string{dir, ".git"}, "/")
	if _, err := os.Stat(dirDotGit); err != nil {
		return
	}
	isGitRepo = true
	return
}

// dir exist && contains .git, do 'git pull'
// dir exist && !contains .git,  remove dir; do 'git clone'
// !dir exist, do 'git clone'
func (pkg PackageReplacer) Merge() error {
	dir, exist, isGitRepo := pkg.ExistOnDisk()
	if exist {
		if isGitRepo {
			cmd := exec.Command("git", "pull", "--progress")
			cmd.Dir = dir
			cmd.Stdout = out
			cmd.Stderr = out
			if err := cmd.Run(); err != nil {
				return errors.WithMessage(err, "do 'git pull'")
			}
			return nil
		}
		if err := os.Remove(dir); err != nil {
			return errors.WithMessage(err, "remove existed dir")
		}
	}
	cmd := exec.Command("git", "clone", "--progress", pkg.GitRemoteRepo, dir)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return errors.WithMessage(err, "do 'git clone'")
	}
	return nil
}
