package mapreduce

import (
	"fmt"
	"sync"
)


//
// schedule() starts and waits for all tasks in the given phase (mapPhase
// or reducePhase). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var nOther int // number of inputs (for reduce) or outputs (for map)

	var wg sync.WaitGroup

	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		nOther = nReduce
	case reducePhase:
		ntasks = nReduce
		nOther = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nOther)

	// All ntasks tasks have to be scheduled on workers. Once all tasks
	// have completed successfully, schedule() should return.
	//
	// Your code here (Part III, Part IV).
	//
	doStr := "Worker.DoTask"
	for i := 0; i < ntasks; i++ {
		var args DoTaskArgs
		args.JobName = jobName
		args.Phase = phase
		args.TaskNumber = i
		args.NumOtherPhase = nOther
		if phase == mapPhase {
			args.File = mapFiles[i]
		} else {
			args.File = ""
		}

		wg.Add(1)
		go func() {
			w := <-registerChan
			call(w, doStr, args, nil)
			select {
			case registerChan <- w:
			default:
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("Schedule: %v done\n", phase)
}
