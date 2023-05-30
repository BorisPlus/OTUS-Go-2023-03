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

	stageOutput := make(Bi)
	for stageID, stage := range stages {

		staged := stage(stageInput)

		go func(stageId int, in In, done In, out Bi) {
			executePipepoint(stageId, in, done, out)
		}(stageID, staged, done, stageOutput)

		fmt.Println("Configute STAGING: stage", stageID)
		fmt.Println("Configute STAGING: stage", stageID, "with done", fmt.Sprintf("%p", done))
		fmt.Println("Configute STAGING: stage", stageID, "with in", fmt.Sprintf("%p", stageInput))
		fmt.Println("Configute STAGING: stage", stageID, "with out", fmt.Sprintf("%p", stageOutput))

		stageInput = stageOutput
		stageOutput = make(Bi)
	}

	fmt.Println("Configute STAGING: out", stageInput)
	fmt.Println()
	return stageInput
}

func executePipepoint(stageID int, in In, done In, out Bi) {
	processor := func(stageId int, in In, out Bi) Out {
		terminated := make(Bi)
		go func() {
			defer func() {
				fmt.Println("stage", stageId, "processor", "end")

				fmt.Println("stage", stageId, "processor", "try to close(out)")
				close(out)
				fmt.Println("stage", stageId, "processor", "close(out)")

				fmt.Println("stage", stageId, "processor", "try setup terminated")
				close(terminated)
				fmt.Println("stage", stageId, "processor", "was setup terminated")
			}()
			for {
				select {
				case value, ok := <-in:
					fmt.Println("stage", stageId, "processor", "get from input", "value", value, "ok", ok)
					if !ok {
						fmt.Println("stage", stageId, "processor", "!ok - return")
						return
					}
					fmt.Println("stage", stageId, "processor", "get from input", "value", value, "ok", ok, "try put to out")
					out <- value
					fmt.Println("stage", stageId, "processor", "get from input", "value", value, "ok", ok, "was put to out")
				case <-done:
					fmt.Println("stage", stageId, "processor", "done - return")
					return
				}
			}
		}()
		return terminated
	}
	terminated := processor(stageID, in, out)
	fmt.Println("stage", stageID, "try", "terminated")
	<-terminated
	fmt.Println("stage", stageID, "was", "terminated")
}
