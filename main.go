package main

import (
	"context"
	"sync"

	"github.com/shikro/osmodel/processor"
	"github.com/shikro/osmodel/scheduler"
	"github.com/shikro/osmodel/taskgenerator"
)

func main() {
	p := processor.New()
	s := scheduler.New()
	p.SetScheduler(s)
	s.SetProcessor(p)

	ctx := context.Background()
	p.Start(ctx)
	s.Start(ctx)

	var wg sync.WaitGroup
	generator := taskgenerator.New(&wg)
	for range 9 {
		s.ScheldueTask(generator.Next(), true)
	}
	wg.Wait()
}
