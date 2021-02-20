package main

import (
	_ "embed" //needed for embedding files
	"math/rand"
	"time"

	"github.com/slurdge/goeland/cmd"
	"github.com/slurdge/goeland/version"
)

//go:embed CHANGELOG.md
var changeLog string

func main() {
	version.ExtractVersionFromChangelog(changeLog)
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
