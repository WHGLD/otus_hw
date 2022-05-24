package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {

	out := worker(in, done)

	for _, stage := range stages {
		if stage != nil {
			out = stage(worker(out, done))
		}
	}

	return out
}

func worker(in In, done In) Out {
	localCh := make(Bi)

	go func() {
		defer close(localCh)

		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				localCh <- value
			}
		}
	}()

	return localCh
}
