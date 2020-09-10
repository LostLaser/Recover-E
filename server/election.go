package server

// Election is an interface for all supported election types
type Election interface {
	startElection(*Server)
	connectServers(map[string]*Server) map[string]*Server
}
