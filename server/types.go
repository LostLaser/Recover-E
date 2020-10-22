package server

import (
	"time"

	"github.com/LostLaser/election/server/communication"
)

// Process describes all classes that can be run as a mock server
type Process interface {
	Boot()
	Restart()
	Stop()
	Print()
	GetID() string
	SetMaster(string) bool
	IsUp() bool
}

// Setup describes all classes that can be used to construct an array of servers
type Setup interface {
	Setup(int, *communication.Emitter, time.Duration) map[string]Process
}
