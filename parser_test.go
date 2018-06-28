package dfon

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testFiles, err := ioutil.ReadDir("./test_data/")
	if assert.NoError(t, err, "error while reading test directory") {
		for i := range testFiles {
			data, err := ioutil.ReadFile("./test_data/" + testFiles[i].Name())
			if assert.NoError(t, err, "error while reading test file", testFiles[i].Name()) {
				_, err := Parse(bytes.NewBuffer(data))
				assert.NoError(t, err, "error while parsing")
			}
		}
	}
}
