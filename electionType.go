package election

// Election is an interface for all supported election types
type Election interface {
	StartElection(*Server)
	ConnectServers(map[string]*Server)
}
