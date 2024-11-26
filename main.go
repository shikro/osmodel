package main

import (
	"context"
	"sync"
	"time"

	"github.com/shikro/osmodel/events"
	"github.com/shikro/osmodel/processor"
	"github.com/shikro/osmodel/scheduler"
	"github.com/shikro/osmodel/taskgenerator"
)

func main() {
	p := processor.New()
	pingEvent := events.NewPingEvent(1 * time.Second)
	s := scheduler.New(pingEvent.Event)
	p.SetScheduler(s)
	s.SetProcessor(p)

	ctx := context.Background()
	pingEvent.Start(ctx)
	p.Start(ctx)
	s.Start(ctx)

	var wg sync.WaitGroup
	generator := taskgenerator.New(&wg)
	for range 9 {
		s.ScheldueTask(generator.Next())
	}
	wg.Wait()
}
