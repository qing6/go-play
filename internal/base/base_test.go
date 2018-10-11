package base

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInTesting(t *testing.T) {
	assert.True(t, InTesting())
}

func TestFields_exportTo(t *testing.T) {
	var fields Fields
	assert.Nil(t, fields)

	buf := new(bytes.Buffer)
	fields.exportTo(buf)
	assert.Equal(t, "", buf.String())

	fields = map[string]interface{}{
		"field1": 111,
		"field2": "222",
	}
	buf.Reset()
	fields.exportTo(buf)
	assert.Equal(t, " field1=111 field2=222", buf.String())
}

func TestSplitFilepath(t *testing.T) {
	var dir string
	var filename string
	var fileExt string

	dir, filename, fileExt = SplitFilepath("/a/b/c/d")
	assert.Equal(t, "/a/b/c", dir)
	assert.Equal(t, "d", filename)
	assert.Equal(t, "", fileExt)
}
