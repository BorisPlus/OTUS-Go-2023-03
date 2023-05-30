package hw06pipelineexecution

import (
	"fmt"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageInput := in
	fmt.Println("Configute STAGING: in", stageInput)

	for stageID, stage := range stages {
		stagedInput := stage(stageInput)
		stageOutput := executePipepoint(stageID, stagedInput, done)

		fmt.Println("Configute STAGING: stage", stageID)
		fmt.Println("Configute STAGING: stage", stageID, "with done", fmt.Sprintf("%p", done))
		fmt.Println("Configute STAGING: stage", stageID, "with in", fmt.Sprintf("%p", stageInput))
		fmt.Println("Configute STAGING: stage", stageID, "with out", fmt.Sprintf("%p", stageOutput))

		stageInput = stageOutput
	}

	fmt.Println("Configute STAGING: out", stageInput)
	fmt.Println()
	return stageInput
}

func executePipepoint(stageID int, in In, done In) Bi {
	out := make(Bi)
	go func() {
		defer func() {
			fmt.Println("stage", stageID, "processor", "end")
			fmt.Println("stage", stageID, "processor", "try to close(out)")
			close(out)
			fmt.Println("stage", stageID, "processor", "close(out)")
		}()
		for {
			select {
			case value, ok := <-in:
				fmt.Println("stage", stageID, "processor", "get from input", "value", value, "ok", ok)
				if !ok {
					fmt.Println("stage", stageID, "processor", "!ok - return")
					return
				}
				fmt.Println("stage", stageID, "processor", "get from input", "value", value, "ok", ok, "try put to out")
				out <- value
				fmt.Println("stage", stageID, "processor", "get from input", "value", value, "ok", ok, "was put to out")
			case <-done:
				fmt.Println("stage", stageID, "processor", "done - return")
				return
			}
		}
	}()
	return out
}
