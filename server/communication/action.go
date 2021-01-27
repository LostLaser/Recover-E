package communication

// Action describes possible interactions that can happen with the election cluster
type Action string

const (
	Heartbeat        Action = "HEARTBEAT"
	Started                 = "STARTED"
	Stopped                 = "STOPPED"
	NotMaster               = "NOT_MASTER"
	ElectionStarted         = "ELECTION_STARTED"
	Elect                   = "ELECT"
	Elected                 = "ELECTED"
	ElectionEnded           = "ELECTION_ENDED"
	StartNewElection        = "START_NEW_ELECTION"
)
