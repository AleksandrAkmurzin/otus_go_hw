package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	in = stageStoppable(in, done)

	for _, stage := range stages {
		in = stageStoppable(stage(in), done)
	}

	return in
}

func stageStoppable(in In, done In) Out {
	inDouble := make(Bi)

	go func() {
		defer func() {
			close(inDouble)
			// Prevent goroutines leak.
			for range in {
			}
		}()

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				default:
					inDouble <- v
				}
			}
		}
	}()

	return inDouble
}
