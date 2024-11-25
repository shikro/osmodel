package processor

import (
	"context"

	"github.com/shikro/osmodel/task"
)

type Scheduler interface {
	ScheldueTask(task task.Task, new bool)
	TaskReady()
}

type Processor struct {
	scheduler Scheduler
	newTask   chan task.Task
	task      task.Task
}

func New() *Processor {
	return &Processor{
		scheduler: nil,
		newTask:   make(chan task.Task),
		task:      nil,
	}
}

func (p *Processor) SetScheduler(s Scheduler) {
	p.scheduler = s
}

func (p *Processor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case newTask := <-p.newTask:
				if p.task != nil && !p.task.Finished() {
					p.scheduler.ScheldueTask(p.task, false)
				}
				p.task = newTask
			default:
				if p.task != nil {
					p.task.DoWork()
					if p.task.Finished() {
						p.task = nil
						p.scheduler.TaskReady()
					}
				}
			}
		}
	}()
}

func (p *Processor) ExecuteTask(task task.Task) {
	p.newTask <- task
}
