package message

// Notify contains a notice to set your current master
type Notify struct {
	Master  string
	visited []string
}

// NewNotify creates a new instance of an Notify message
func NewNotify(id string) Notify {
	e := Notify{}
	e.visited = make([]string, 0)
	e.Master = id
	return e
}

// Visited returns whether or not the id is in the visited list
func (e *Notify) Visited(id string) bool {
	for _, v := range e.visited {
		if v == id {
			return true
		}
	}

	return false
}

// AddVisited appends the input server id to the visited list
func (e *Notify) AddVisited(id string) {
	e.visited = append(e.visited, id)
}
