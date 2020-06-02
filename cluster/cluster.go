package cluster

import (
	"fmt"

	"github.com/LostLaser/recover-e/emitter"
	"github.com/LostLaser/recover-e/server"
)

// Cluster is a linked collection of servers
type Cluster struct {
	linkedServers map[string]*server.Server
	emitter       *emitter.Emitter
}

// New will create a cluster with the specified number of servers
func New(serverCount int) *Cluster {
	c := new(Cluster)
	c.linkedServers = make(map[string]*server.Server)
	c.emitter = emitter.New(100)

	for i := 0; i < serverCount; i++ {
		s := server.New(c.emitter)
		c.linkedServers[s.GetID()] = s
	}
	for currKey, currserver := range c.linkedServers {
		for key, server := range c.linkedServers {
			if currKey != key {
				currserver.NeighborServers[key] = server
			}
		}
		go currserver.Start()
	}

	return c
}

//ListServers prints all servers in the cluster
func (c Cluster) ListServers() {
	for _, s := range c.linkedServers {
		s.Print()
	}
}

// Purge will stop all of the linked servers
func (c Cluster) Purge() {
	for _, s := range c.linkedServers {
		s.Stop()
	}
}

//Stream will retrieve an ordered event log of the servers' actions
func (c Cluster) Stream() {
	for {
		msg := c.emitter.Read()
		fmt.Println(msg["from"], "sent", msg["details"], "to", msg["to"])
	}
}
