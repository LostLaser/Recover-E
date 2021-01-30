package election

import (
	"errors"
	"time"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/communication"
)

// Cluster is a linked collection of servers
type Cluster struct {
	linkedServers map[string]server.Process
	emitter       *communication.Emitter
}

// New will create a cluster with the specified number of servers
func New(processSetup server.Setup, serverCount int, heartbeatPause time.Duration) *Cluster {
	c := new(Cluster)
	c.linkedServers = make(map[string]server.Process)
	c.emitter = communication.New(serverCount * 10)

	for k, v := range processSetup.Setup(serverCount, c.emitter, heartbeatPause) {
		c.linkedServers[k] = v
	}

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
func (c Cluster) ReadEvent() interface{} {
	return c.emitter.Read()
}

func (c Cluster) getServerByID(id string) (server.Process, error) {
	for key, s := range c.linkedServers {
		if id == key {
			return s, nil
		}
	}
	return nil, errors.New("No server found with specified ID")
}
