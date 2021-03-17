package server

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneration(t *testing.T) {
	i := GenerateUniqueID()
	assert.Regexp(t, regexp.MustCompile("^[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}$"), i)
}
