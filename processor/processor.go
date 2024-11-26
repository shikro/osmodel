package processor

import (
	"context"

	"github.com/shikro/osmodel/task"
)

type Scheduler interface {
	RetakeTask(task task.Task)
	TaskDone()
	TaskWaiting(task task.Task)
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
					p.scheduler.RetakeTask(p.task)
				}
				p.task = newTask
			default:
				if p.task != nil {
					p.task.DoWork()
					if p.task.Finished() {
						p.task = nil
						p.scheduler.TaskDone()
					} else if extTask, ok := p.task.(task.ExtendedTask); ok {
						if extTask.Wait() {
							p.task = nil
							p.scheduler.TaskWaiting(extTask)
						}
					}
				}
			}
		}
	}()
}

func (p *Processor) ExecuteTask(task task.Task) {
	p.newTask <- task
}
