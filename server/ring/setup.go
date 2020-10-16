package ring

import (
	"time"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/communication"
)

// Setup is used in the cluster to create an array of servers of the ring type
type Setup struct {
}

// Setup links together the specified amount of bully election servers together
func (s Setup) Setup(c int, e *communication.Emitter, t time.Duration) map[string]server.Process {
	lp := make(map[string]server.Process)

	sa := make([]*Process, c)

	for i := 0; i < c; i++ {
		sa[i] = New(e, t)
	}

	for i := 0; i < c; i++ {
		for k := (i + 1) % c; k != i; k = (k + 1) % c {
			sa[i].linkedServers = append(sa[i].linkedServers, sa[k])
		}
		lp[sa[i].GetID()] = sa[i]
	}

	return lp
}
