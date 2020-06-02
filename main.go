package main

import (
	"time"

	"github.com/LostLaser/recover-e/cluster"
)

func main() {
	c1 := cluster.New(10)
	go c1.Stream()
	time.Sleep(time.Second * 4)
	c1.Purge()

	return
}
