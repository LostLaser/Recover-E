package server

import (
	"testing"
	"time"

	"github.com/LostLaser/election/server/communication"
	"github.com/stretchr/testify/assert"
)

var heartbeatPause = time.Millisecond * 250

func TestNew(t *testing.T) {
	b := New(communication.New(10), heartbeatPause)
	assert.False(t, b.IsUp())
	assert.Equal(t, heartbeatPause, b.HeartbeatPause)
}

func TestRestart(t *testing.T) {
	b := New(communication.New(10), heartbeatPause)
	b.Restart()
	assert.True(t, b.IsUp())
}

func TestStop(t *testing.T) {
	b := New(communication.New(10), heartbeatPause)
	b.Restart()
	b.Stop()
	assert.False(t, b.IsUp())
	assert.Empty(t, b.Master)
}

func TestPrint(t *testing.T) {
	return // TODO
}

func TestGetID(t *testing.T) {
	b := New(communication.New(10), heartbeatPause)
	assert.Equal(t, b.GetID(), b.ID)
}

func TestSetMasterUp(t *testing.T) {
	masterID := "1234"
	b := New(communication.New(10), heartbeatPause)
	b.Restart()
	assert.True(t, b.SetMaster(masterID))
	assert.Equal(t, masterID, b.Master)
}

func TestSetMasterDown(t *testing.T) {
	masterID := "1234"
	b := New(communication.New(10), heartbeatPause)
	assert.False(t, b.SetMaster(masterID))
	assert.Equal(t, "", b.Master)
}

func TestIsUp(t *testing.T) {
	b := New(communication.New(10), heartbeatPause)
	assert.False(t, b.IsUp())
	b.Restart()
	assert.True(t, b.IsUp())
}
