package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func sendToPipeline(done In, terminate Bi, in In, out Bi) {
	for data := range in {
		select {
		case <-done:
			close(out)
			close(terminate)
			return
		case out <- data:
		}
	}
	close(out)
}

func runStage(done In, in In, out Bi, stage Stage) {
	for data := range stage(in) {
		select {
		case <-done:
			close(out)
			return
		case out <- data:
		}
	}
	close(out)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	terminate := make(Bi)
	toStage := make(Bi)
	fromStage := make(Bi)

	go sendToPipeline(done, terminate, in, toStage)
	for i, stage := range stages {
		stage := stage
		if i > 0 {
			toStage = fromStage
			fromStage = make(Bi)
		}
		go runStage(terminate, toStage, fromStage, stage)
	}
	return fromStage
}
