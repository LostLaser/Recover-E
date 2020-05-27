package server

import (
	"time"
)

func (s *Server) run() {
	for s.state == running {
		time.Sleep(time.Second)
		if !s.pingMaster() {
			startElection(s)
		}
	}
}

func (s *Server) pingMaster() bool {
	if s.master == "" || (s.master != s.id && !s.NeighborServers[s.master].isUp()) {
		return false
	}
	return true
}

func (s *Server) setMaster(masterID string) {
	s.master = masterID
}

func (s *Server) isUp() bool {
	return s.state == running
}
