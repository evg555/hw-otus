package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = executeStage(in, done, stage)
	}

	return in
}

func executeStage(in In, done In, stage Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		if in == nil {
			return
		}

		stageOut := stage(in)

		for {
			select {
			case <-done:
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}

				out <- v
			}
		}
	}()

	return out
}
