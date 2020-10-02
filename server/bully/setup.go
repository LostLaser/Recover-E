package bully

import (
	"time"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/communication"
)

// Setup is used in the cluster to create an array of servers of the bully type
type Setup struct {
}

// Setup links together and creates the specified amount of bully election servers together
func (s Setup) Setup(c int, e *communication.Emitter, t time.Duration) map[string]server.Process {
	lp := make(map[string]server.Process)
	sa := make([]*Process, c)

	for i := 0; i < c; i++ {
		sa[i] = New(e, t)
	}

	for i := 0; i < c; i++ {
		for k := 0; k < c; k++ {
			if sa[i].GetID() != sa[k].GetID() {
				sa[i].NeighborServers[sa[k].GetID()] = sa[k]
			}
		}
		lp[sa[i].GetID()] = sa[i]
	}

	return lp
}
