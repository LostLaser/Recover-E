package election

import (
	"errors"
	"time"
)

// Cluster is a linked collection of servers
type Cluster struct {
	linkedServers     map[string]*Server
	emitter           *Emitter
	electionAlgorithm Election
}

// New will create a cluster with the specified number of servers
func New(serverCount int, heartbeatPause time.Duration, algorithm Election) *Cluster {
	c := new(Cluster)
	c.linkedServers = make(map[string]*Server)
	c.emitter = NewEmitter(serverCount * 10)
	c.electionAlgorithm = algorithm

	for i := 0; i < serverCount; i++ {
		s := NewServer(c.emitter, heartbeatPause, c.electionAlgorithm)
		c.linkedServers[s.GetID()] = s
	}

	c.electionAlgorithm.ConnectServers(c.linkedServers)

	for _, currserver := range c.linkedServers {
		go currserver.Initialize()
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
		s.Start()
	}
	return err
}

//ReadEvent will retrieve a single event log of the servers' actions
func (c Cluster) ReadEvent() map[string]string {
	return c.emitter.Read()
}

func (c Cluster) getServerByID(id string) (*Server, error) {
	for key, s := range c.linkedServers {
		if id == key {
			return s, nil
		}
	}
	return nil, errors.New("No server found with specified ID")
}
