package election

import (
	"errors"
	"time"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/communication"
)

// Cluster is a linked collection of servers
type Cluster struct {
	linkedServers     map[string]*server.Server
	emitter           *communication.Emitter
	electionAlgorithm server.Algorithm
}

// New will create a cluster with the specified number of servers
func New(serverCount int, heartbeatPause time.Duration, algorithm server.Algorithm) *Cluster {
	c := new(Cluster)
	c.linkedServers = make(map[string]*server.Server)
	c.emitter = communication.New(serverCount * 10)
	c.electionAlgorithm = algorithm

	for i := 0; i < serverCount; i++ {
		s := server.New(c.emitter, heartbeatPause, c.electionAlgorithm)
		c.linkedServers[s.GetID()] = s
	}

	c.electionAlgorithm.ConnectServers(c.linkedServers)

	for _, currserver := range c.linkedServers {
		go currserver.Boot()
	}

	return c
}

//ServerIds returns all server ids in the cluster
func (c Cluster) ServerIds() []string {
	var ids []string
	for _, s := range c.linkedServers {
		ids = append(ids, s.GetID())
	}

	return ids
}

// Purge will stop all of the linked servers
func (c Cluster) Purge() {
	for _, s := range c.linkedServers {
		s.Stop()
	}
}

// StopServer stops the server with the specified id in the cluster
func (c Cluster) StopServer(id string) error {
	s, err := c.getServerByID(id)
	if err == nil {
		s.Stop()
	}
	return err
}

// StartServer starts the server with the specified id in the cluster
func (c Cluster) StartServer(id string) error {
	s, err := c.getServerByID(id)
	if err == nil {
		s.Restart()
	}
	return err
}

//ReadEvent will retrieve a single event log of the servers' actions
func (c Cluster) ReadEvent() map[string]string {
	return c.emitter.Read()
}

func (c Cluster) getServerByID(id string) (*server.Server, error) {
	for key, s := range c.linkedServers {
		if id == key {
			return s, nil
		}
	}
	return nil, errors.New("No server found with specified ID")
}
