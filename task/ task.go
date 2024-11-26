package task

import (
	"fmt"
	"sync"
	"time"
)

type Task interface {
	ID() string
	DoWork()
	Finished() bool
	Priority() int
}

type ExtendedTask interface {
	Task
	Wait() bool
}

type SimpleTask struct {
	id       string
	size     uint
	progress uint
	priority int
	wg       *sync.WaitGroup
}

func NewSimpleTask(id string, size uint, priority int, wg *sync.WaitGroup) *SimpleTask {
	return &SimpleTask{
		id:       id,
		size:     size,
		progress: 0,
		priority: priority,
		wg:       wg,
	}
}

func (t SimpleTask) ID() string {
	return t.id
}

func (t *SimpleTask) DoWork() {
	if t.progress < t.size {
		time.Sleep(100 * time.Millisecond)
		t.progress++
		fmt.Printf("%s | progress: %d/%d | priority: %d\n", t.id, t.progress, t.size, t.priority)
		if t.Finished() {
			t.wg.Done()
		}
	}
}

func (t SimpleTask) Finished() bool {
	return t.progress == t.size
}

func (t SimpleTask) Priority() int {
	return t.priority
}

type TaskWithWait struct {
	SimpleTask
	waited bool
}

func NewTaskWithWait(id string, size uint, priority int, wg *sync.WaitGroup) *TaskWithWait {
	return &TaskWithWait{
		SimpleTask: *NewSimpleTask(id, size, priority, wg),
		waited:     false,
	}
}

func (t *TaskWithWait) Wait() bool {
	if !t.waited && t.progress >= t.size/2 {
		fmt.Printf("%s start waiting\n", t.id)
		t.waited = true
		return true
	}
	return false
}
