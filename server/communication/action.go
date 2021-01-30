package communication

// Action describes possible interactions that can happen with the election cluster
type Action string

const (
	// Heartbeat indicates a check up message between a follower and master server
	Heartbeat Action = "HEARTBEAT"
	// Started indicates a server has successfully started
	Started = "STARTED"
	// Stopped indicates a server has gone inactive
	Stopped = "STOPPED"
	// NotMaster should be emit by a server when their reign has ended as master
	NotMaster = "NOT_MASTER"
	// ElectionStarted indicates a server has gone into election processing
	ElectionStarted = "ELECTION_STARTED"
	// Elect indicates a decision by one server that its' leader has been declared
	Elect = "ELECT"
	// Elected indicates a server has been named the leader
	Elected = "ELECTED"
	// ElectionEnded is the signal that a server has stopped its' election algorithm
	ElectionEnded = "ELECTION_ENDED"
	// StartNewElection is used by some election algorithms to notify another server that
	// it should start its' election process
	StartNewElection = "START_NEW_ELECTION"
)
