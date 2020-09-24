package server

import "time"

func (s *Server) run() {
	for {
		if s.state == running {
			if !s.pingMaster() || s.triggerElection {
				s.electionAlgorithm.StartElection(s)
				s.triggerElection = false
			}
		}
		time.Sleep(s.heartbeatPause)
	}
}

func (s *Server) pingMaster() bool {
	s.emitter.Write(s.id, s.master, "HEARTBEAT")
	if s.master == "" || (s.master != s.id && !s.NeighborServers[s.master].isUp()) {
		return false
	}

	return true
}
