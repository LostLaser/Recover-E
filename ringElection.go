package election

// RingElection implements the ring election algorithm
type RingElection struct {
}

// StartElection runs the ring election algorithm on the provided server
func (r *RingElection) StartElection(s *Server) {
	return
}

// ConnectServers links the input servers in accordance with the ring election algorithm
func (r *RingElection) ConnectServers(s map[string]*Server) {
	for currKey, currserver := range s {
		for key, server := range s {
			if currKey != key {
				currserver.NeighborServers[key] = server
			}
		}
	}
}
