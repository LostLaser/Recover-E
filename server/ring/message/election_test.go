package message

import (
	"testing"

	"github.com/LostLaser/election/server"
	"github.com/stretchr/testify/assert"
)

func TestNewElection(t *testing.T) {
	id := server.GenerateUniqueID()
	e := NewElection(id)

	assert.True(t, e.Exists(id), "Message not added to default notification list")
}

func TestAddNotified(t *testing.T) {
	id := server.GenerateUniqueID()
	addID := server.GenerateUniqueID()
	e := NewElection(id)
	e.AddNotified(addID)

	found := false
	for _, v := range e.notified {
		if v == addID {
			found = true
			break
		}
	}
	assert.True(t, found, "Message not added to notification list")
}

func TestNotifiedCount(t *testing.T) {
	id := server.GenerateUniqueID()
	addID := server.GenerateUniqueID()
	e := NewElection(id)
	e.AddNotified(addID)

	assert.Equal(t, 2, e.NotifiedCount(), "Message not added to notification list")
}

func TestGetHighest(t *testing.T) {
	id := server.GenerateUniqueID()
	addID := server.GenerateUniqueID()
	e := NewElection(id)
	e.AddNotified(addID)

	highest := ""
	if id > addID {
		highest = id
	} else {
		highest = addID
	}

	assert.Equal(t, highest, e.GetHighest(), "Highest id not determined correctly")
}

func TestExistsOk(t *testing.T) {
	id := server.GenerateUniqueID()
	addID := server.GenerateUniqueID()
	e := NewElection(id)
	e.AddNotified(addID)

	assert.True(t, e.Exists(addID), "Existant ID was not found")
}

func TestExistsNotFound(t *testing.T) {
	id := server.GenerateUniqueID()
	e := NewElection(id)

	assert.False(t, e.Exists(server.GenerateUniqueID()), "Non-existant ID was found")
}
