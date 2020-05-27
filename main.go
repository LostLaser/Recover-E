package main

import (
	"fmt"
	"time"

	"github.com/LostLaser/recover-e/cluster"
)

func main() {
	fmt.Println("HEY there!")
	c1 := cluster.New(10)
	time.Sleep(time.Second * 4)
	c1.ListServers()
	c1.Purge()

	return
}
