package communication

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControlJson(t *testing.T) {
	expectedTarget := "1234"
	expectedAction := Started

	c := NewControl(expectedTarget, Action(expectedAction))

	json, err := json.Marshal(c)
	if err != nil {
		assert.Fail(t, "Json could not be parsed for Control struct", err)
	}

	expectedLabels := []string{
		"target",
		"action",
		"type",
	}

	for _, label := range expectedLabels {
		assert.Contains(t, string(json), label)
	}
}

func TestEventJson(t *testing.T) {
	expectedFrom := "1234"
	expectedTo := "5678"
	expectedAction := Heartbeat

	c := NewEvent(expectedFrom, expectedTo, Action(expectedAction))

	json, err := json.Marshal(c)
	if err != nil {
		assert.Fail(t, "Json could not be parsed for Event struct", err)
	}

	expectedLabels := []string{
		"from",
		"to",
		"action",
		"type",
	}

	for _, label := range expectedLabels {
		assert.Contains(t, string(json), label+"\"")
	}
}
