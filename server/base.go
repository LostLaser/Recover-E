package server

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/LostLaser/election/server/communication"
)

// Base is a single entity
type Base struct {
	Master         string
	ID             string
	State          int
	ElectionLock   sync.Mutex
	Emitter        *communication.Emitter
	HeartbeatPause time.Duration
}

// New will create a cluster with the specified number of servers
func New(e *communication.Emitter, heartbeatPause time.Duration) *Base {
	s := new(Base)
	s.ID = GenerateUniqueID()
	s.State = Stopped
	s.Emitter = e
	s.HeartbeatPause = heartbeatPause

	return s
}

// Restart the provided server
func (s *Base) Restart() {
	s.State = Running
	s.Emitter.Write(communication.NewControl(s.ID, communication.Started))
}

// Stop the provided server
func (s *Base) Stop() {
	s.State = Stopped
	s.Master = ""
	s.Emitter.Write(communication.NewControl(s.ID, communication.Stopped))
}

// String displays the server information in a readable format
func (s *Base) String() string {
	str := fmt.Sprint("ID:", s.ID, "Master:", s.Master)
	return str
}

// GetID returns the server id
func (s *Base) GetID() string {
	return s.ID
}

// SetMaster assigns the specified master to the calling server's master variable
func (s *Base) SetMaster(masterID string) bool {
	if !s.IsUp() {
		return false
	}
	s.ElectionLock.Lock()
	defer s.ElectionLock.Unlock()
	if masterID != s.ID && s.ID == s.Master {
		s.Emitter.Write(communication.NewControl(s.ID, communication.NotMaster))
	}
	s.Master = masterID

	return true
}

// IsUp returns whether or not the current server is running
func (s *Base) IsUp() bool {
	return s.State == Running
}

// MarshalJSON retrieves the target server as a json string
func (s *Base) MarshalJSON() ([]byte, error) {
	r := struct {
		ID             string
		HeartbeatPause time.Duration
		Master         string
		State          int
	}{
		s.ID,
		s.HeartbeatPause,
		s.Master,
		s.State,
	}

	return json.Marshal(r)
}
