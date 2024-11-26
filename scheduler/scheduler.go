package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/shikro/osmodel/task"
)

type Processor interface {
	ExecuteTask(task task.Task)
}

const (
	PrioritiesCount = 4
	maxTasks        = 100
)

type Scheduler struct {
	processor       Processor
	event           func() chan struct{}
	curTaskPriority int
	tasksCountMu    sync.RWMutex
	tasksCount      uint
	queues          [][]task.Task
	waitingQueue    []task.Task
	newTask         chan task.Task
	retakeTask      chan task.Task
	taskDone        chan struct{}
	taskWaiting     chan task.Task
}

func New(e func() chan struct{}) *Scheduler {
	return &Scheduler{
		processor:       nil,
		event:           e,
		curTaskPriority: PrioritiesCount,
		tasksCount:      0,
		queues:          make([][]task.Task, PrioritiesCount),
		waitingQueue:    make([]task.Task, 0),
		newTask:         make(chan task.Task),
		retakeTask:      make(chan task.Task, 1),
		taskDone:        make(chan struct{}, 1),
		taskWaiting:     make(chan task.Task, 1),
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
			case task := <-s.taskWaiting:
				s.curTaskPriority = PrioritiesCount
				s.waitingQueue = append(s.waitingQueue, task)
				t := s.selectTask()
				if t != nil {
					s.curTaskPriority = t.Priority()
					s.processor.ExecuteTask(t)
				}
			case <-s.taskDone:
				s.curTaskPriority = PrioritiesCount
				t := s.selectTask()
				if t != nil {
					s.curTaskPriority = t.Priority()
					s.processor.ExecuteTask(t)
				}
			case <-s.event():
				if len(s.waitingQueue) > 0 {
					task := s.waitingQueue[0]
					s.waitingQueue = s.waitingQueue[1:]
					p := task.Priority()
					s.queues[p] = append(s.queues[p], nil)
					copy(s.queues[p][1:], s.queues[p])
					s.queues[p][0] = task
					fmt.Printf("%s awaken\n", task.ID())

					t := s.selectTask()
					if t != nil {
						s.curTaskPriority = t.Priority()
						s.processor.ExecuteTask(t)
					}
				}
			}
		}
	}()
}

func (s *Scheduler) ScheldueTask(t task.Task) error {
	s.tasksCountMu.Lock()
	if s.tasksCount >= maxTasks {
		return fmt.Errorf("max task count reached")
	}
	s.tasksCount++
	s.tasksCountMu.Unlock()

	s.newTask <- t
	return nil
}

func (s *Scheduler) RetakeTask(t task.Task) {
	s.retakeTask <- t
}

func (s *Scheduler) TaskDone() {
	s.taskDone <- struct{}{}
}

func (s *Scheduler) TaskWaiting(t task.Task) {
	s.taskWaiting <- t
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
