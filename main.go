package main

import (
	"math/rand"
	"time"

	"github.com/slurdge/goeland/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()

}
