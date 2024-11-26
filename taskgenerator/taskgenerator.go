package taskgenerator

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/shikro/osmodel/task"
)

type TaskGenerator struct {
	counter int
	wg      *sync.WaitGroup
}

func New(wg *sync.WaitGroup) *TaskGenerator {
	return &TaskGenerator{
		wg: wg,
	}
}

func (t *TaskGenerator) Next() task.Task {
	t.wg.Add(1)
	t.counter++
	if rand.IntN(100) < 50 {
		return task.NewTaskWithWait(fmt.Sprintf("ext-task-%d", t.counter), 5, rand.IntN(4), t.wg)
	}
	return task.NewSimpleTask(fmt.Sprintf("smp-task-%d", t.counter), 5, rand.IntN(4), t.wg)
}
