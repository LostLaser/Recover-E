package election

//Emitter is a specialized queue for messaging
type Emitter struct {
	messages chan map[string]string
}

//NewEmitter creates an instance of an emitter
func NewEmitter(bufferSize int) *Emitter {
	e := new(Emitter)
	e.messages = make(chan map[string]string, bufferSize)

	return e
}

//Write will add a new message to the emitter
func (e *Emitter) Write(from string, to string, action string) {
	e.messages <- map[string]string{"from": from, "to": to, "action": action}
}

//Read returns the oldest message in the emitter, will block if no message is available
func (e *Emitter) Read() map[string]string {
	val := <-e.messages
	return val
}
