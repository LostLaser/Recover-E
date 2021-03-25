package ring

import (
	"regexp"
	"testing"
	"time"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/communication"
	"github.com/stretchr/testify/assert"
)

var pollingTime = time.Millisecond * 250
var settleTime = pollingTime + time.Millisecond*100

func TestNew(t *testing.T) {
	r := New(communication.New(10), pollingTime)

	assert.False(t, r.IsUp())
	assert.Regexp(t, regexp.MustCompile("^[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}$"), r.ID)
	assert.Equal(t, r.HeartbeatPause, pollingTime)
}

func TestBootNoLinkedServers(t *testing.T) {
	expectedState := server.Running
	r := New(communication.New(10), pollingTime)

	go r.Boot()
	time.Sleep(settleTime)

	assert.Equal(t, expectedState, r.State)
	// TODO: assert.Equal(t, r.ID, r.Master, "Lone server not elected as master")
}

func TestBootLinkedServers(t *testing.T) {
	r := New(communication.New(10), pollingTime)
	r2 := New(communication.New(10), pollingTime)

	r.linkedServers = append(r.linkedServers, r2)
	r2.linkedServers = append(r2.linkedServers, r)
	go r.Boot()
	go r2.Boot()
	time.Sleep(settleTime)

	masterID := ""
	if r.ID > r2.ID {
		masterID = r.ID
	} else if r.ID < r2.ID {
		masterID = r2.ID
	} else {
		assert.FailNow(t, "UUIDs are equivalent. Consider different ID gen solution or buy a lottery ticket")
	}
	assert.Equal(t, masterID, r.Master, "Server with greatest ID not elected master")
	assert.Equal(t, masterID, r2.Master, "Server with greatest ID not elected master")
}

func TestStoppedMaster(t *testing.T) {
	r1 := New(communication.New(10), pollingTime)
	r2 := New(communication.New(10), pollingTime)
	r3 := New(communication.New(10), pollingTime)

	r1.linkedServers = append(r1.linkedServers, r2, r3)
	r2.linkedServers = append(r2.linkedServers, r3, r1)
	r3.linkedServers = append(r3.linkedServers, r1, r2)
	servers := map[string]*Process{r1.ID: r1, r2.ID: r2, r3.ID: r3}
	go r1.Boot()
	go r2.Boot()
	go r3.Boot()
	time.Sleep(settleTime)

	// set highest ID as expected master
	expectedMaster := ""
	for _, s := range servers {
		if s.ID > expectedMaster {
			expectedMaster = s.ID
		}
	}
	for _, s := range servers {
		assert.Equal(t, expectedMaster, s.Master, "Master not equal to max ID")
	}

	servers[expectedMaster].Stop()
	time.Sleep(settleTime)

	// set highest active ID as expected master
	expectedMasterAfter := ""
	for _, s := range servers {
		if s.ID > expectedMasterAfter && s.IsUp() {
			expectedMasterAfter = s.ID
		}
	}
	for _, s := range servers {
		if s.IsUp() {
			assert.Equal(t, expectedMasterAfter, s.Master, "Master not equal to max ID")
		}
	}

}
