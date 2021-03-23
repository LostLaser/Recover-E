package election

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/communication"
	"go.uber.org/zap"
)

// Cluster is a linked collection of servers
type Cluster struct {
	linkedServers map[string]server.Process
	emitter       *communication.Emitter
	logger        *zap.SugaredLogger
	ID            string
}

// New will create a cluster with the specified number of servers
func New(processSetup server.Setup, serverCount int, heartbeatPause time.Duration, logger *zap.Logger) *Cluster {
	c := new(Cluster)
	c.ID = server.GenerateUniqueID()
	c.linkedServers = make(map[string]server.Process)
	c.emitter = communication.New(serverCount * 10)
	c.logger = logger.Sugar()

	for k, v := range processSetup.Setup(serverCount, c.emitter, heartbeatPause) {
		c.linkedServers[k] = v
	}

	for _, currserver := range c.linkedServers {
		go currserver.Boot()
	}
	c.logger.Debugf("Created cluster: %v", c)
	return c
}

//ServerIds returns all server ids in the cluster
func (c Cluster) ServerIds() []string {
	var ids []string
	for _, s := range c.linkedServers {
		ids = append(ids, s.GetID())
	}

	return ids
}

// Purge will stop all of the linked servers but not clear their references
func (c Cluster) Purge() {
	c.logger.Debugf("Stopping all servers in cluster %s", c.ID)
	for _, s := range c.linkedServers {
		s.Stop()
	}
}

// StopServer stops the server with the specified id in the cluster
func (c Cluster) StopServer(id string) error {
	c.logger.Debugf("Attempting to stop server %s", id)
	s, err := c.getServerByID(id)
	if err == nil {
		s.Stop()
	} else {
		c.logger.Infof("%v", err)
	}
	return err
}

// StartServer starts the server with the specified id in the cluster
func (c Cluster) StartServer(id string) error {
	c.logger.Debugf("Attempting to start server %s", id)
	s, err := c.getServerByID(id)
	if err == nil {
		s.Restart()
	} else {
		c.logger.Infof("%v", err)
	}
	return err
}

//ReadEvent will retrieve a single event log of the servers' actions
func (c Cluster) ReadEvent() interface{} {
	ev := c.emitter.Read()
	c.logger.Debugf("Emitter message: %v", ev)
	return ev
}

// MarshalJSON retrieves the target cluster as a json string
func (c Cluster) MarshalJSON() ([]byte, error) {
	var serverArray []server.Process
	for _, s := range c.linkedServers {
		serverArray = append(serverArray, s)
	}

	r := struct {
		ID          string
		ServerCount int
		Servers     []server.Process
	}{
		c.ID,
		len(c.linkedServers),
		serverArray,
	}

	return json.Marshal(r)

}

// String gets the string representation of the target cluster
func (c Cluster) String() string {
	v, err := c.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%s", v)
}

func (c Cluster) getServerByID(id string) (server.Process, error) {
	if s, found := c.linkedServers[id]; found {
		return s, nil
	}

	return nil, errors.New("No server found with ID '" + id + "'")
}
