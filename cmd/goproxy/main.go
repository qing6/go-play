// Fork
package main

import (
	"flag"
	"os"
	"path/filepath"
	"qing/cmds/internal/base"
	"qing/cmds/pkg/model"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	listen string
	gopath string
	logger *logrus.Logger
)

func init() {
	flag.StringVar(&listen, "listen", "0.0.0.0:9066", "service listen address")
	flag.Parse()

	gopath = filepath.SplitList(os.Getenv("GOPATH"))[0]
	if len(gopath) == 0 {
		logger.Panic("Not Found $GOPATH")
	}

	logger = base.GetLogger()
}

func main() {
}

// func parsePseudoVersion() (*gitCommitSelector, error) {
// 	return nil, errors.New("not impl!")
// }

type ModuleExporter interface {
	Export(model.PackageReplacer, PeseudoVersion)
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
