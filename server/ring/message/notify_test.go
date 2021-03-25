package message

import (
	"testing"

	"github.com/LostLaser/election/server"
	"github.com/stretchr/testify/assert"
)

func TestNewNotify(t *testing.T) {
	id := server.GenerateUniqueID()
	n := NewNotify(id)
	assert.Equal(t, id, n.Master, "Base master ID not equal to input ID")
}

func TestAddVisited(t *testing.T) {
	id := server.GenerateUniqueID()
	addID := server.GenerateUniqueID()
	n := NewNotify(id)
	n.AddVisited(addID)

	assert.Equal(t, addID, n.visited[0], "ID not added to notify object")
}

func TestVisitedFound(t *testing.T) {
	id := server.GenerateUniqueID()
	addID := server.GenerateUniqueID()
	n := NewNotify(id)
	n.AddVisited(addID)

	assert.True(t, n.Visited(addID), "ID not found in visited list")
}

func TestVisitedNotFound(t *testing.T) {
	id := server.GenerateUniqueID()
	addID := server.GenerateUniqueID()
	n := NewNotify(id)

	assert.False(t, n.Visited(addID), "ID found in visited list")
}
