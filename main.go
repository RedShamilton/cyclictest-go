// Author: Sean Hamilton <skhamilt@eng.ucsd.edu>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package main

import (
	"os/user"
	"fmt"
	"os"
	"flag"
	"sync"
	"syscall"
	"time"
	"github.com/RedShamilton/cyclictest-go/types"
)

var wg sync.WaitGroup
var running = true
var params []types.TaskParameters

func main() {
	// Check to ensure this is being run as root
	u,_ := user.Current()
	if u.Uid != "0" {
		fmt.Fprintln(os.Stderr, "cyclictest must be run as root")
		os.Exit(int(syscall.EPERM))
	}

	// Set flag parsing to printout usage and exit on error
	//flag.CommandLine.Init("", flag.ExitOnError)

	var numLoops uint
	flag.UintVar(&numLoops, "l", 1000, "number of `loops`")
	flag.UintVar(&numLoops, "loops", 1000, "")

	var priority int
	flag.IntVar(&priority, "p", 0, "priority for highest `priority` thread")
	flag.IntVar(&priority, "priority", 0, "")

	var numTasks uint
	flag.UintVar(&numTasks, "t", 1, "number of `tasks`")
	flag.UintVar(&numTasks, "tasks", 1, "")

	var distance time.Duration
	flag.DurationVar(&distance,"d", 500*time.Microsecond, "`distance` of task intervals")
	flag.DurationVar(&distance,"distance", 500*time.Microsecond, "")

	var interval time.Duration
	flag.DurationVar(&interval,"i", 1000*time.Microsecond, "base `interval` of task")
	flag.DurationVar(&interval,"interval", 1000*time.Microsecond, "")

	flag.Parse()
	//TODO: check for valid ranges of flags

        wg.Add(int(numTasks))
	nextInterval := interval
	nextPrio := priority
	params = make([]types.TaskParameters, numTasks)
	for i := uint(0); i < numTasks; i++ {
		params[i].Init(i,nextInterval,int32(nextPrio),new(types.TaskStatistics))
		nextPrio -= 1
		nextInterval += distance
		go worker(&params[i], numLoops)
	}

        wg.Wait()

		for i := uint(0); i < numTasks; i++ {
			params[i].Stats.PrintResults()
			time.Sleep(10*time.Microsecond)
		}
}