package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPseudoVersion1(t *testing.T) {
	var v *pseudoVersion1
	var err error

	s := "123"
	v, err = NewPseudoVersion1(s)
	assert.Nil(t, v)
	assert.Equal(t, "regex match fail.", err.Error())

	s = "v01.0.0-01234567890123-012345678901"
	v, err = NewPseudoVersion1(s)
	assert.Nil(t, v)
	assert.Equal(t, "regex match fail.", err.Error())

	s = "v23000000000000000000000000000000000000000000.0.0-01234567890123-012345678901"
	v, err = NewPseudoVersion1(s)
	assert.Nil(t, v)
	assert.True(t, strings.HasPrefix(err.Error(),
		"get version major fail. s=23000000000000000000000000000000000000000000"))

	s = "v0.0.0-01234567890123-012345678901"
	v, err = NewPseudoVersion1(s)
	assert.Nil(t, v)
	assert.True(t, strings.HasPrefix(err.Error(),
		"get commit time fail. s=01234567890123"))
	t.Log(err.Error())

	s = "v0.0.0-20181011010203-0123456789af"
	v, err = NewPseudoVersion1(s)
	assert.Nil(t, err)
	assert.Equal(t, pseudoVersion1{
		version:          [3]int{0, 0, 0},
		commitAt:         time.Date(2018, time.October, 11, 01, 02, 03, 0, time.UTC),
		commitHashPrefix: "0123456789af",
	}, *v)
}
