package bully

import (
	"regexp"
	"testing"
	"time"

	"github.com/LostLaser/election/server/communication"
	"github.com/stretchr/testify/assert"
)

var pollingTime time.Duration = time.Millisecond * 250
var settleTime time.Duration = pollingTime + time.Millisecond*100

func TestNew(t *testing.T) {

	b := New(communication.New(10), pollingTime)
	assert.False(t, b.IsUp())
	assert.Regexp(t, regexp.MustCompile("^[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}$"), b.ID)
	assert.Equal(t, b.HeartbeatPause, pollingTime)
}

func TestBootNoNeighbors(t *testing.T) {
	b := New(communication.New(10), pollingTime)
	go b.Boot()
	time.Sleep(time.Millisecond * 50)
	assert.True(t, b.IsUp())
	time.Sleep(pollingTime + settleTime)
	assert.True(t, b.pingMaster())
}

func TestPingMaster(t *testing.T) {
	b := New(communication.New(10), pollingTime)
	assert.False(t, b.pingMaster())
}

func TestIsHighest(t *testing.T) {
	b := New(communication.New(10), pollingTime)
	b2 := New(communication.New(10), pollingTime)
	b.NeighborServers[b2.ID] = b2
	b2.NeighborServers[b.ID] = b

	go b.Boot()
	go b2.Boot()
	time.Sleep(settleTime)

	if b.ID > b2.ID {
		assert.True(t, b.isHighest())
	} else {
		assert.True(t, b2.isHighest())
	}
}

func TestStartElection(t *testing.T) {
	b := New(communication.New(10), pollingTime)
	b2 := New(communication.New(10), pollingTime)
	b.NeighborServers[b2.ID] = b2
	b2.NeighborServers[b.ID] = b

	// Boot should eventually trigger startElection after a settle time
	go b.Boot()
	go b2.Boot()
	time.Sleep(settleTime)

	if b.ID > b2.ID {
		assert.True(t, b.isHighest())
		assert.True(t, b2.Master == b.ID)
		assert.True(t, b.Master == b.ID)
	} else {
		assert.True(t, b2.isHighest())
		assert.True(t, b.Master == b2.ID)
		assert.True(t, b2.Master == b2.ID)
	}
}
