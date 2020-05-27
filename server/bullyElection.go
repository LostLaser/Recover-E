package server

func startElection(s *Server) {
	if isHighest(s) {
		notifyLow(s)
		s.setMaster(s.id)
	}
}

func isHighest(s *Server) bool {
	for key, neighbor := range s.NeighborServers {
		if key > s.id && neighbor.isUp() {
			return false
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
