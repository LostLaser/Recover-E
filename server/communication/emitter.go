package communication

// Emitter is a specialized queue for messaging
type Emitter struct {
	messages chan interface{}
}

// New creates an instance of an emitter
func New(bufferSize int) *Emitter {
	e := new(Emitter)
	e.messages = make(chan interface{}, bufferSize)

	return e
}

// Write will add a new message to the emitter
func (e *Emitter) Write(m interface{}) {
	e.messages <- m
}

// Read returns the oldest message in the emitter, will block if no message is available
func (e *Emitter) Read() interface{} {
	val := <-e.messages
	return val
}
