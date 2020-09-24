package server

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/LostLaser/election/server/communication"
)

// Server is a single entity
type Server struct {
	master            string
	id                string
	state             string
	NeighborServers   map[string]*Server
	electionAlgorithm Algorithm
	electionLock      sync.Mutex
	triggerElection   bool
	emitter           *communication.Emitter
	heartbeatPause    time.Duration
}

const (
	running = "running"
	stopped = "stopped"
)

// New will create a cluster with the specified number of servers
func New(e *communication.Emitter, heartbeatPause time.Duration, electionAlgorithm Algorithm) *Server {
	s := new(Server)
	s.id = generateUniqueID()
	s.state = running
	s.NeighborServers = make(map[string]*Server)
	s.emitter = e
	s.heartbeatPause = heartbeatPause
	s.electionAlgorithm = electionAlgorithm

	return s
}

// Boot brings up the server and runs main process
func (s *Server) Boot() {
	s.state = running
	s.run()
}

// Restart the provided server
func (s *Server) Restart() {
	s.state = running
	s.emitter.Write(s.id, "", "STARTED")
}

// Stop the provided server
func (s *Server) Stop() {
	s.state = stopped
	s.master = ""
	s.emitter.Write(s.id, "", "STOPPED")
}

// Print displays the server information in a readable format
func (s *Server) Print() {
	fmt.Println("ID:", s.id, " Master:", s.master)
}

// GetID returns the server id
func (s *Server) GetID() string {
	return s.id
}

func generateUniqueID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

func (s *Server) setMaster(masterID string) {
	if !s.isUp() {
		return
	}
	s.electionLock.Lock()
	defer s.electionLock.Unlock()
	if masterID != s.id && s.id == s.master {
		s.emitter.Write(s.id, "", "NOT_MASTER")
	}
	s.emitter.Write(masterID, s.id, "ELECT")
	s.master = masterID
}

func (s *Server) isUp() bool {
	return s.state == running
}
