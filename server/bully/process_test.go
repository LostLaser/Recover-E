package bully

import (
	"regexp"
	"testing"
	"time"

	"github.com/LostLaser/election/server/communication"
	"github.com/stretchr/testify/assert"
)

var settleTime time.Duration = time.Millisecond * 250

func TestNew(t *testing.T) {
	e := communication.New(10)
	interval := time.Millisecond

	b := New(e, interval)
	assert.False(t, b.IsUp())
	assert.Regexp(t, regexp.MustCompile("^[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}$"), b.ID)
	assert.Equal(t, b.HeartbeatPause, interval)
}

func TestBootNoNeighbors(t *testing.T) {
	interval := time.Second

	b := New(communication.New(10), interval)
	go b.Boot()
	assert.True(t, b.IsUp())
	time.Sleep(interval + settleTime)
	assert.True(t, b.pingMaster())
}

func TestPingMaster(t *testing.T) {
	b := New(communication.New(10), time.Millisecond)
	assert.False(t, b.pingMaster())
}
