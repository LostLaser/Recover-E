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
	if actualServerCount != expectedServerCount {
		t.Errorf("Server count was incorrect, got: %d, want: %d.", actualServerCount, expectedServerCount)
	}
}

func TestNeighbors(t *testing.T) {

}

func TestServerListingCount(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)

	actualServerCount := len(cluster.ServerIds())
	if actualServerCount != expectedServerCount {
		t.Errorf("Server count was incorrect, got: %d, want: %d.", actualServerCount, expectedServerCount)
	}
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
		if !found {
			t.Errorf("Didn't find linked server %s in server listing", i)
		}
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
		t.Error("No event received after expected time interval")
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

	if len(cluster.ServerIds()) == 0 {
		t.Error("Test requires at least one server in cluster")
		return
	}
	id := serverIds[0]

	err := cluster.StopServer(id)
	if err != nil {
		t.Errorf("Error recieved for invalid id: %s, error: %s", id, err)
	}
}

func TestStopInvl(t *testing.T) {
	expectedServerCount := 3
	id := "invl"
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)
	err := cluster.StopServer("invl")
	if err == nil {
		t.Errorf("No error recieved for invalid id: %s", id)
	}
}

func TestStart(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)
	serverIds := cluster.ServerIds()

	if len(cluster.ServerIds()) == 0 {
		t.Error("Test requires at least one server in cluster")
		return
	}
	id := serverIds[0]

	err := cluster.StartServer(id)
	if err != nil {
		assert.FailNow(t, "Could not start the server")
	}
}

func TestStartInvl(t *testing.T) {
	expectedServerCount := 3
	id := "invl"
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime, logger)
	err := cluster.StartServer(id)
	if err == nil {
		assert.FailNow(t, "No error recieved for invalid id")
	}
}

func TestMarshalJSON(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}
	cluster := New(setup, expectedServerCount, cycleTime, logger)

	_, err := cluster.MarshalJSON()
	if err != nil {
		assert.FailNow(t, "Unexpected error when marshalling to json", err)
	}
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
