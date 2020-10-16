package message

// Elected contains a notice to set your current master
type Elected struct {
	Master  string
	visited []string
}

// NewElected creates a new instance of an Elected message
func NewElected(id string) Elected {
	e := Elected{}
	e.Master = id
	return e
}

// Visited returns whether or not the id is in the visited list
func (e Elected) Visited(id string) bool {
	for _, v := range e.visited {
		if v == id {
			return true
		}
	}

	return false
}

// AddVisited appends the input server id to the visited list
func (e Elected) AddVisited(id string) {
	e.visited = append(e.visited, id)
}
