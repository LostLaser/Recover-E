package cluster

import (
	"github.com/LostLaser/recover-e/server"
)

// Cluster is a linked collection of servers
type Cluster struct {
	linkedServers map[string]*server.Server
}

// New will create a cluster with the specified number of servers
func New(serverCount int) *Cluster {
	c := new(Cluster)
	c.linkedServers = make(map[string]*server.Server)

	for i := 0; i < serverCount; i++ {
		s := server.New()
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
