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
	notifyQueue   chan message.Notify
}

// New will create a cluster with the specified number of servers
func New(e *communication.Emitter, heartbeatPause time.Duration) *Process {
	b := new(Process)
	b.ID = server.GenerateUniqueID()
	b.State = server.Stopped
	b.Emitter = e
	b.HeartbeatPause = heartbeatPause
	b.electionQueue = make(chan message.Election, 4)
	b.notifyQueue = make(chan message.Notify, 4)

	return b
}

// Boot brings up the server and runs main process
func (r *Process) Boot() {
	r.State = server.Running
	go r.electionResponder()
	r.run()
}

func (r *Process) run() {
	for {
		time.Sleep(r.HeartbeatPause)
		if r.State == server.Running && !r.pingMaster() {
			r.startElection()
		}
	}
}

func (r *Process) electionResponder() {
	for {
		select {
		case m := <-r.electionQueue:
			// add name to alive list and send to neighbor
			if m.Exists(r.ID) {
				r.SetMaster(m.GetHighest())
				q := message.NewNotify(m.GetHighest())
				q.AddVisited(r.ID)
				r.getNeighbor().notifyQueue <- q
			} else {
				r.Emitter.Write(communication.NewEvent(r.ID, r.getNeighbor().ID, communication.StartNewElection))
				m.AddNotified(r.ID)
				r.getNeighbor().electionQueue <- m
			}
		case m := <-r.notifyQueue:
			// set master to consensus and send to neighbor
			if !m.Visited(r.ID) {
				if r.ID == m.Master {
					r.Emitter.Write(communication.NewControl(r.ID, communication.Elected))
				}
				r.SetMaster(m.Master)
				m.AddVisited(r.ID)
				r.Emitter.Write(communication.NewEvent(r.ID, r.getNeighbor().ID, communication.Elect))
				r.getNeighbor().notifyQueue <- m
			}
		}
	}
}

func (r *Process) startElection() {
	p := r.getNeighbor()
	if p == nil {
		// only server in cluster
		r.SetMaster(r.ID)
		return
	}
	r.Emitter.Write(communication.NewEvent(r.ID, r.getNeighbor().ID, communication.StartNewElection))
	p.electionQueue <- message.NewElection(r.ID)
}

func (r *Process) pingMaster() bool {
	r.Emitter.Write(communication.NewEvent(r.ID, r.Master, communication.Heartbeat))
	if r.Master == "" || (r.Master != r.ID && !r.getServer(r.Master).IsUp()) {
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
