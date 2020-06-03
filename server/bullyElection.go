package server

type bullyElection struct {
}

func startElection(s *Server) {
	s.emitter.Write(s.id, "", "ELECTION STARTED")
	if isHighest(s) {
		notifyLow(s)
		s.setMaster(s.id)
	}
	s.emitter.Write(s.id, "", "ELECTION ENDED")
}

func isHighest(s *Server) bool {
	for key, neighbor := range s.NeighborServers {
		if key > s.id {
			s.emitter.Write(s.id, key, "CHECK HIGHER ELECTION")
			if neighbor.isUp() {
				return false
			}
		}
	}
	return true
}

func notifyLow(s *Server) {
	for key, neighbor := range s.NeighborServers {
		if key < s.id {
			neighbor.setMaster(s.id)
		}
	}
}
