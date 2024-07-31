package runner

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func example() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	r := New(bytes.NewBuffer([]byte{}))
	r.Runfp("obvious_nonsense %s %s", "bad", "args")
	return nil
}

func TestHandle(t *testing.T) {
	err := example()
	assert.NotNil(t, err)
}
