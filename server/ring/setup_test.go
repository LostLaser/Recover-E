package ring

import (
	"testing"

	"github.com/LostLaser/election/server/communication"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	e := communication.New(10)
	count := 3
	s := Setup{}
	smap := s.Setup(count, e, pollingTime)
	assert.Equal(t, count, len(smap))
}
