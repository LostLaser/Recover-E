package election

import (
	"testing"
	"time"

	"github.com/LostLaser/election/server/bully"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func TestNew(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)

	actualServerCount := len(cluster.linkedServers)
	assert.Equal(t, expectedServerCount, actualServerCount, "Server count was incorrect")
}

func TestServerListingCount(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)

	actualServerCount := len(cluster.ServerIds())
	assert.Equal(t, expectedServerCount, actualServerCount, "Server count was incorrect")
}

func TestServerListingConsistency(t *testing.T) {
	serverCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, serverCount, cycleTime, logger)

	for _, i := range cluster.ServerIds() {
		found := false
		for k := range cluster.linkedServers {
			if i == k {
				found = true
				break
			}
		}
		assert.True(t, found, "Didn't find linked server %s in server listing")
	}
}

func TestReadEvent(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)
	c := make(chan (int))

	go func() {
		cluster.ReadEvent()
		c <- 1
	}()

	time.Sleep(cycleTime + time.Second/4)
	select {
	case <-c:
		return
	default:
		assert.FailNow(t, "No event received after expected time interval")
	}

}

func TestPurge(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)

	cluster.Purge()
	for _, v := range cluster.linkedServers {
		assert.False(t, v.IsUp(), "One of the servers was not stopped")
	}
}

func TestStop(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)
	serverIds := cluster.ServerIds()

	assert.NotEqual(t, 0, len(cluster.ServerIds()), "Test requires at least one server in cluster")
	id := serverIds[0]

	err := cluster.StopServer(id)
	assert.Nil(t, err, "Error recieved for invalid id")
}

func TestStopInvl(t *testing.T) {
	expectedServerCount := 3
	id := "invl"
	cycleTime := time.Second
	setup := bully.Setup{}
	cluster := New(setup, expectedServerCount, cycleTime, logger)

	err := cluster.StopServer(id)
	assert.NotNil(t, err, "No error recieved for invalid id")

}

func TestStart(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)
	serverIds := cluster.ServerIds()

	assert.NotEqual(t, 0, len(cluster.ServerIds()), "Test requires at least one server in cluster")

	id := serverIds[0]

	err := cluster.StartServer(id)
	assert.Nil(t, err, "Could not start the server", err)
}

func TestStartInvl(t *testing.T) {
	expectedServerCount := 3
	id := "invl"
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)
	err := cluster.StartServer(id)
	assert.NotNil(t, err, "No error recieved for invalid id")

}

func TestMarshalJSON(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}
	cluster := New(setup, expectedServerCount, cycleTime, logger)

	_, err := cluster.MarshalJSON()
	assert.Nil(t, err, "Unexpected error when marshalling to json", err)
}

func TestString(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}
	cluster := New(setup, expectedServerCount, cycleTime, logger)
	cluster.Purge()

	str := cluster.String()

	assert.Contains(t, str, cluster.ID, "String cluster representation does not contain cluster ID")
}
