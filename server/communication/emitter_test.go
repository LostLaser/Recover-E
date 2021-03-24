package communication

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmitterFull(t *testing.T) {
	e := New(1)
	expected := "My message"

	e.Write(expected)
	actual := e.Read()

	assert.Equal(t, expected, actual, "Message exchanged was modified")
}
