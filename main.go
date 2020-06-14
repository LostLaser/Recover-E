package main

import (
	"time"

	"github.com/LostLaser/recoverE/cluster"
)

func main() {
	c := cluster.New(4, time.Second)
	go c.Stream()
	time.Sleep(4 * time.Second)
	c.ListServers()
}
