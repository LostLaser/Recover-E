package election

import (
	"testing"
	"time"

	"github.com/LostLaser/election/server/bully"
)

func TestNew(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime)

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

	cluster := New(setup, expectedServerCount, cycleTime)

	actualServerCount := len(cluster.ServerIds())
	if actualServerCount != expectedServerCount {
		t.Errorf("Server count was incorrect, got: %d, want: %d.", actualServerCount, expectedServerCount)
	}
}

func TestServerListingConsistency(t *testing.T) {
	serverCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, serverCount, cycleTime)

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

	cluster := New(setup, expectedServerCount, cycleTime)
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

func TestStop(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	setup := bully.Setup{}

	cluster := New(setup, expectedServerCount, cycleTime)
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

	cluster := New(setup, expectedServerCount, cycleTime)
	err := cluster.StopServer("invl")
	if err == nil {
		t.Errorf("No error recieved for invalid id: %s", id)
	}
}
