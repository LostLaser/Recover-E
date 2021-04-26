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

// New will create a new instance of a bully server
func New(e *communication.Emitter, heartbeatPause time.Duration) *Process {
	b := new(Process)
	b.ID = server.GenerateUniqueID()
	b.State = server.Stopped
	b.Emitter = e
	b.HeartbeatPause = heartbeatPause
	b.NeighborServers = make(map[string]*Process)

	return b
}

// Boot brings up the server and runs main process. This function blocks while the process is active
func (b *Process) Boot() {
	b.State = server.Running
	b.run()
}

func (b *Process) run() {
	for {
		time.Sleep(b.HeartbeatPause)
		if b.State == server.Running {
			if !b.pingMaster() || b.triggerElection {
				b.startElection()
				b.triggerElection = false
			}
		}
	}
}

func (b *Process) startElection() {
	b.Emitter.Write(communication.NewControl(b.ID, communication.ElectionStarted))
	if b.isHighest() {
		b.notifyLow()
		b.Emitter.Write(communication.NewEvent(b.ID, b.ID, communication.Elect))
		b.SetMaster(b.ID)
		b.Emitter.Write(communication.NewControl(b.ID, communication.Elected))
	}
	b.Emitter.Write(communication.NewControl(b.ID, communication.ElectionEnded))
}

func (b *Process) pingMaster() bool {
	b.Emitter.Write(communication.NewEvent(b.ID, b.Master, communication.Heartbeat))
	if b.Master == "" || (b.Master != b.ID && !b.NeighborServers[b.Master].IsUp()) {
		return false
	}

	return true
}

func (b *Process) isHighest() bool {
	for id, neighbor := range b.NeighborServers {
		if id > b.ID {
			if neighbor.IsUp() {
				b.Emitter.Write(communication.NewEvent(b.ID, id, communication.StartNewElection))
				neighbor.triggerElection = true
				return false
			}
		}
	}
	return true
}

func (b *Process) notifyLow() {
	for key, neighbor := range b.NeighborServers {
		if key < b.ID && neighbor.SetMaster(b.ID) {
			b.Emitter.Write(communication.NewEvent(b.ID, neighbor.ID, communication.Elect))
		}
	}
}
