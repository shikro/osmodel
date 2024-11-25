package task

import (
	"fmt"
	"sync"
	"time"
)

type Task interface {
	DoWork()
	Finished() bool
	Priority() int
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

func (t *SimpleTask) DoWork() {
	if t.progress < t.size {
		t.progress++
		fmt.Printf("%s | progress: %d/%d | priority: %d\n", t.id, t.progress, t.size, t.priority)
		time.Sleep(100 * time.Millisecond)
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

type ExtendedTask struct {
	size     uint
	progress uint
}

func NewExtendedTask(size uint) ExtendedTask {
	return ExtendedTask{size: size, progress: 0}
}

func (t *ExtendedTask) DoWork() {
	if t.size < t.progress {
		t.progress++
	}
}

func (t ExtendedTask) Finished() bool {
	return t.progress == t.size
}
