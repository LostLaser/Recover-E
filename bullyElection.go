package election

// BullyElection implements the bully election algorithm
type BullyElection struct {
}

// StartElection runs the bully election algorithm on the provided server
func (b BullyElection) StartElection(s *Server) {
	s.emitter.Write(s.id, "", "ELECTION_STARTED")
	if isHighest(s) {
		notifyLow(s)
		s.setMaster(s.id)
		s.emitter.Write(s.id, "", "ELECTED")
	}
	s.emitter.Write(s.id, "", "ELECTION_ENDED")
}

// ConnectServers links the input servers in accordance with the bully algorithm
func (b BullyElection) ConnectServers(s map[string]*Server) {
	for currKey, currserver := range s {
		for key, server := range s {
			if currKey != key {
				currserver.NeighborServers[key] = server
			}
		}
	}
}

func isHighest(s *Server) bool {
	for id, neighbor := range s.NeighborServers {
		if id > s.id {
			if neighbor.isUp() {
				s.emitter.Write(s.id, id, "START_NEW_ELECTION")
				neighbor.triggerElection = true
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
