package server

const (
	// Running is the state that server will be set to when online
	Running = iota
	// Stopped is the state that server will be set to when offline
	Stopped = iota
)
