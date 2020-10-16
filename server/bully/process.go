package bully

import (
	"time"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/communication"
)

// Process extends Base to implement the bully algorithm
type Process struct {
	server.Base
	NeighborServers map[string]*Process
	triggerElection bool
}

// New will create a cluster with the specified number of servers
func New(e *communication.Emitter, heartbeatPause time.Duration) *Process {
	b := new(Process)
	b.ID = server.GenerateUniqueID()
	b.State = server.Running
	b.Emitter = e
	b.HeartbeatPause = heartbeatPause
	b.NeighborServers = make(map[string]*Process)

	return b
}

// Boot brings up the server and runs main process
func (b *Process) Boot() {
	b.State = server.Running
	b.run()
}

func (b *Process) run() {
	for {
		if b.State == server.Running {
			if !b.pingMaster() || b.triggerElection {
				b.startElection()
				b.triggerElection = false
			}
		}
		time.Sleep(b.HeartbeatPause)
	}
}

func (b *Process) startElection() {
	b.Emitter.Write(b.ID, "", "ELECTION_STARTED")
	if b.isHighest() {
		b.notifyLow()
		b.SetMaster(b.ID)
		b.Emitter.Write(b.ID, "", "ELECTED")
	}
	b.Emitter.Write(b.ID, "", "ELECTION_ENDED")
}

func (b *Process) pingMaster() bool {
	b.Emitter.Write(b.ID, b.Master, "HEARTBEAT")
	if b.Master == "" || (b.Master != b.ID && !b.NeighborServers[b.Master].IsUp()) {
		return false
	}

	return true
}

func (b *Process) isHighest() bool {
	for id, neighbor := range b.NeighborServers {
		if id > b.ID {
			if neighbor.IsUp() {
				b.Emitter.Write(b.ID, id, "START_NEW_ELECTION")
				neighbor.triggerElection = true
				return false
			}
		}
	}
	return true
}

func (b *Process) notifyLow() {
	for key, neighbor := range b.NeighborServers {
		if key < b.ID {
			neighbor.SetMaster(b.ID)
		}
	}
}
