package ring

import "github.com/LostLaser/election/server"

// Process implements the ring election algorithm
type Process struct {
	server.Base
	NeighborServers map[string]*Process
	triggerElection bool
}

// Boot brings up the server and runs main process
func (r *Process) Boot() {
	r.State = server.Running
}

// ConnectServers links the input servers in accordance with the bully algorithm
func (r *Process) ConnectServers(s map[string]*server.Base) {
	return
}
