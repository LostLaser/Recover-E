package server

// Algorithm is an interface for all supported election types
type Algorithm interface {
	StartElection(*Server)
	ConnectServers(map[string]*Server)
}
