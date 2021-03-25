package election

import (
	"testing"
	"time"

	"github.com/LostLaser/election/server/bully"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var defaultServerCount = 3
var defaultCycleTime = time.Second

func defaultTestCluster(t *testing.T) *Cluster {
	setup := bully.Setup{}
	logger := zap.NewNop()

	return New(setup, defaultServerCount, defaultCycleTime, logger)
}

func TestNew(t *testing.T) {
	cluster := defaultTestCluster(t)

	actualServerCount := len(cluster.linkedServers)
	assert.Equal(t, defaultServerCount, actualServerCount, "Server count was incorrect")
}

func TestServerListingCount(t *testing.T) {
	cluster := defaultTestCluster(t)

	actualServerCount := len(cluster.ServerIds())
	assert.Equal(t, defaultServerCount, actualServerCount, "Server count was incorrect")
}

func TestServerListingConsistency(t *testing.T) {
	cluster := defaultTestCluster(t)

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
	cluster := defaultTestCluster(t)
	c := make(chan (int))

	go func() {
		cluster.ReadEvent()
		c <- 1
	}()

	time.Sleep(defaultCycleTime + time.Second/4)
	select {
	case <-c:
		return
	default:
		assert.FailNow(t, "No event received after expected time interval")
	}

}

func TestPurge(t *testing.T) {
	cluster := defaultTestCluster(t)

	cluster.Purge()
	for _, v := range cluster.linkedServers {
		assert.False(t, v.IsUp(), "One of the servers was not stopped")
	}
}

func TestStop(t *testing.T) {
	cluster := defaultTestCluster(t)
	serverIds := cluster.ServerIds()

	assert.NotEqual(t, 0, len(cluster.ServerIds()), "Test requires at least one server in cluster")
	id := serverIds[0]

	err := cluster.StopServer(id)
	assert.Nil(t, err, "Error recieved for invalid id")
}

func TestStopInvl(t *testing.T) {
	id := "invl"

	cluster := defaultTestCluster(t)

	err := cluster.StopServer(id)
	assert.NotNil(t, err, "No error recieved for invalid id")

}

func TestStart(t *testing.T) {
	cluster := defaultTestCluster(t)
	serverIds := cluster.ServerIds()

	assert.NotEqual(t, 0, len(cluster.ServerIds()), "Test requires at least one server in cluster")

	id := serverIds[0]

	err := cluster.StartServer(id)
	assert.Nil(t, err, "Could not start the server")
}

func TestStartInvl(t *testing.T) {
	id := "invl"
	cluster := defaultTestCluster(t)

	err := cluster.StartServer(id)

	assert.NotNil(t, err, "No error recieved for invalid id")

}

func TestMarshalJSON(t *testing.T) {
	cluster := defaultTestCluster(t)

	_, err := cluster.MarshalJSON()
	assert.Nil(t, err, "Unexpected error when marshalling to json")
}

func TestString(t *testing.T) {
	cluster := defaultTestCluster(t)
	cluster.Purge()

	str := cluster.String()

	assert.Contains(t, str, cluster.ID, "Cluster string representation does not contain cluster ID")
}
