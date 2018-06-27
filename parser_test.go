package dfon

import (
	"bytes"
	"io/ioutil"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data, err := ioutil.ReadFile("./test_data/entity_default.txt")
	if assert.NoError(t, err) {
		head, err := Parse(bytes.NewBuffer(data))
		if assert.NoError(t, err) {
			head.Print(os.Stdout)
		}
	}
}
