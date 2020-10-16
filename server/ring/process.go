package ring

import (
	"time"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/communication"
	"github.com/LostLaser/election/server/ring/message"
)

// Process implements the ring election algorithm
type Process struct {
	server.Base
	linkedServers []*Process
	electionQueue chan message.Election
	electedQueue  chan message.Elected
}

// New will create a cluster with the specified number of servers
func New(e *communication.Emitter, heartbeatPause time.Duration) *Process {
	b := new(Process)
	b.ID = server.GenerateUniqueID()
	b.State = server.Running
	b.Emitter = e
	b.HeartbeatPause = heartbeatPause

	return b
}

// Boot brings up the server and runs main process
func (r *Process) Boot() {
	r.State = server.Running
	r.run()
	go r.electionResponder()
}

func (r *Process) run() {
	for {
		if r.State == server.Running && !r.pingMaster() {
			r.startElection()
		}
		time.Sleep(r.HeartbeatPause)
	}
}

func (r *Process) electionResponder() {
	for {
		select {
		case m := <-r.electionQueue:
			// add name to alive list and send to neighbor
			if m.Exists(r.ID) {
				r.SetMaster(m.GetHighest())
				r.getNeighbor().electedQueue <- message.NewElected(m.GetHighest())
			} else {
				m.AddNotified(r.ID)
				r.getNeighbor().electionQueue <- m
			}
		case m := <-r.electedQueue:
			// set master to consensus and send to neighbor
			r.SetMaster(m.Master)
			if !m.Visited(r.ID) {
				m.AddVisited(r.ID)
				r.getNeighbor().electedQueue <- m
			}
		}
	}
}

func (r *Process) startElection() {
	p := r.getNeighbor()
	if p == nil {
		return
	}
	p.electionQueue <- message.NewElection(r.ID)
}

func (r *Process) pingMaster() bool {
	r.Emitter.Write(r.ID, r.Master, "HEARTBEAT")
	if r.getServer(r.Master) == nil || (r.Master != r.ID && !r.getServer(r.Master).IsUp()) {
		return false
	}

	return true
}

func (r *Process) getNeighbor() *Process {
	for _, p := range r.linkedServers {
		if p.IsUp() {
			return p
		}
	}

	return nil
}

func (r *Process) getServer(id string) *Process {
	for _, p := range r.linkedServers {
		if p.ID == id {
			return p
		}
	}

	return nil
}
