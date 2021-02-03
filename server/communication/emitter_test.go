package communication

import "testing"

func TestEmitterFull(t *testing.T) {
	e := New(1)
	expected := "My message"

	e.Write(expected)
	if actual := e.Read(); actual != expected {
		t.Errorf("Message was incorrect got: %s, want: %s.", actual, expected)
	}
}
