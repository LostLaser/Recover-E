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

	return lp
}
