// Fork
package main

import (
	"flag"
	"time"

	"github.com/pkg/errors"
)

var listen string
var modules string

func init() {
	flag.StringVar(&listen, "listen", "0.0.0.0:9066", "service listen address")
	flag.Parse()
}

func main() {
}

// https://tip.golang.org/cmd/go/#hdr-Pseudo_versions
// TODO 这个方法并不完善, 只是为了解析已知的情况:
//   golang.org/x/crypto@v0.0.0-20180904163835-0709b304e793
// 也就是第一种
func parsePseudoVersion() (*gitCommitSelector, error) {
	return nil, errors.New("not impl!")
}

type peseudoVersion interface {
}

// vX.0.0-yyyymmddhhmmss-abcdefabcdef
// ^v(\d+).0.0-(\d{14})-([0-9|a-f]{6})$
type pseudoVersion1 struct {
	version          [3]int // major.minor.patch, default = -1.-1.-1
	commitAt         time.Time
	commitHashPrefix string
}

// TODO add regex
func NewPseudoVersion1(s string) (*pseudoVersion1, error) {
	return nil, errors.New("not impl!")
}
