package main

import (
	"math/rand"
	"time"

	"github.com/slurdge/goeland/cmd"
)

//go:generate go-bindata -prefix "asset/" -pkg cmd -o cmd/bindata.go asset/

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()

}
