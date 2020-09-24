package election

import (
	"testing"
	"time"

	"github.com/LostLaser/election/server"
)

func TestNew(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	algorithm := &server.BullyElection{}

	cluster := New(expectedServerCount, cycleTime, algorithm)

	actualServerCount := len(cluster.linkedServers)
	if actualServerCount != expectedServerCount {
		t.Errorf("Server count was incorrect, got: %d, want: %d.", actualServerCount, expectedServerCount)
	}
}

func TestNeighbors(t *testing.T) {
	expectedServerCount := 5
	expectedNeighborCount := expectedServerCount - 1
	cycleTime := time.Second
	algorithm := &server.BullyElection{}

	cluster := New(expectedServerCount, cycleTime, algorithm)
	for _, server := range cluster.linkedServers {
		neighborCount := len(server.NeighborServers)
		if neighborCount != expectedNeighborCount {
			t.Errorf("Server neighbor count was incorrect, got: %d, want: %d.", neighborCount, expectedNeighborCount)
		}
	}
}

func TestServerListingCount(t *testing.T) {
	expectedServerCount := 3
	cycleTime := time.Second
	algorithm := &server.BullyElection{}

	cluster := New(expectedServerCount, cycleTime, algorithm)

	actualServerCount := len(cluster.ServerIds())
	if actualServerCount != expectedServerCount {
		t.Errorf("Server count was incorrect, got: %d, want: %d.", actualServerCount, expectedServerCount)
	}
}

func TestServerListingConsistency(t *testing.T) {
	serverCount := 3
	cycleTime := time.Second
	algorithm := &server.BullyElection{}

	cluster := New(serverCount, cycleTime, algorithm)

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
	algorithm := &server.BullyElection{}

	cluster := New(expectedServerCount, cycleTime, algorithm)
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
	algorithm := &server.BullyElection{}

	cluster := New(expectedServerCount, cycleTime, algorithm)
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
	cycleTime := time.Second
	algorithm := &server.BullyElection{}
	id := "invl"

	cluster := New(expectedServerCount, cycleTime, algorithm)
	err := cluster.StopServer("invl")
	if err == nil {
		t.Errorf("No error recieved for invalid id: %s", id)
	}
}
