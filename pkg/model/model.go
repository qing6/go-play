package model

import (
	"os"
	"os/exec"
	"path/filepath"
	"qing/cmds/internal/base"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry
var GOPATH string

func init() {
	logger = base.GetLogger().WithField("package", "pkg/model")
	GOPATH = filepath.SplitList(os.Getenv("GOPATH"))[0]
	if len(GOPATH) == 0 {
		logger.Panic("Not Found $GOPATH")
	}
}

type PackageReplacer struct {
	ImportPath    string
	GitRemoteRepo string
}

func (pkg PackageReplacer) InGOPATH() string {
	return filepath.Join(GOPATH, "src", pkg.ImportPath)
}

func (pkg PackageReplacer) ExistInGOPATH() (dir string, exist bool, isGitRepo bool) {
	dir = pkg.InGOPATH()
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
func (pkg PackageReplacer) MergeInGOPATH() error {
	dir, exist, isGitRepo := pkg.ExistInGOPATH()
	if exist {
		if isGitRepo {
			cmd := exec.Command("git", "pull", "--progress")
			cmd.Dir = dir
			cmd.Stdout = logger.Logger.Out
			cmd.Stderr = logger.Logger.Out
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
	cmd.Stdout = logger.Logger.Out
	cmd.Stderr = logger.Logger.Out
	if err := cmd.Run(); err != nil {
		return errors.WithMessage(err, "do 'git clone'")
	}
	return nil
}

// https://tip.golang.org/cmd/go/#hdr-Pseudo_versions
// 根据已出现的情况, 实现第一种, 也就是 type pseudoVersion1 struct
type PeseudoVersion interface {
}

// vX.0.0-yyyymmddhhmmss-abcdefabcdef
// 粗略的正则: ^v(0|[1-9]\d+).0.0-(\d{14})-([0-9|a-f]{12})$
type pseudoVersion1 struct {
	version          [3]int // major.minor.patch
	commitAt         time.Time
	commitHashPrefix string
}

func NewPseudoVersion1(s string) (*pseudoVersion1, error) {
	results := regexp.MustCompile("^v(0|[1-9]\\d+).0.0-(\\d{14})-([0-9|a-f]{12})$").
		FindAllStringSubmatch(s, -1)
	if results == nil {
		return nil, base.Fail("regex match fail.", nil)
	}
	logger.Debug("regex for PseudoVersion1 match results: ", results)

	version := [3]int{0, 0, 0}

	majorNum, err := strconv.Atoi(results[0][1])
	if err != nil {
		return nil, base.FailBy(err, "get version major fail.", base.Fields{"s": results[0][1]})
	}
	version[0] = int(majorNum)

	commitAt, err := time.Parse("20060102150405", results[0][2])
	if err != nil {
		return nil, base.FailBy(err, "get commit time fail.", base.Fields{"s": results[0][2]})
	}

	return &pseudoVersion1{
		version:          version,
		commitAt:         commitAt,
		commitHashPrefix: results[0][3],
	}, nil
}
