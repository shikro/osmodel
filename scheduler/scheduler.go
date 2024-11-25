package scheduler

import (
	"context"

	"github.com/shikro/osmodel/task"
)

type Processor interface {
	ExecuteTask(task task.Task)
}

const (
	prioritiesCount = 4
)

type Scheduler struct {
	processor       Processor
	curTaskPriority int
	queues          [][]task.Task
	newTask         chan task.Task
	retakeTask      chan task.Task
	needTask        chan struct{}
}

func New() *Scheduler {
	return &Scheduler{
		processor:       nil,
		curTaskPriority: prioritiesCount,
		queues:          make([][]task.Task, prioritiesCount),
		newTask:         make(chan task.Task),
		retakeTask:      make(chan task.Task),
		needTask:        make(chan struct{}),
	}
}

func (s *Scheduler) SetProcessor(p Processor) {
	s.processor = p
}

func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case task := <-s.newTask:
				s.queues[task.Priority()] = append(s.queues[task.Priority()], task)
				t := s.selectTask()
				if t != nil {
					s.curTaskPriority = t.Priority()
					s.processor.ExecuteTask(t)
				}
			case task := <-s.retakeTask:
				p := task.Priority()
				s.queues[p] = append(s.queues[p], nil)
				copy(s.queues[p][1:], s.queues[p])
				s.queues[p][0] = task
			case <-s.needTask:
				s.curTaskPriority = prioritiesCount
				t := s.selectTask()
				if t != nil {
					s.curTaskPriority = t.Priority()
					s.processor.ExecuteTask(t)
				}
			}
		}
	}()
}

func (s *Scheduler) ScheldueTask(task task.Task, new bool) {
	if new {
		s.newTask <- task
	} else {
		s.retakeTask <- task
	}
}

func (s *Scheduler) TaskReady() {
	s.needTask <- struct{}{}
}

func (s *Scheduler) selectTask() task.Task {
	var task task.Task
	for i, q := range s.queues {
		if s.curTaskPriority > i && len(q) > 0 {
			task = q[0]
			s.queues[i] = s.queues[i][1:]
			break
		}
	}
	return task
}
