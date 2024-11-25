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

func (t *TaskGenerator) Next() *task.SimpleTask {
	t.wg.Add(1)
	t.counter++
	return task.NewSimpleTask(fmt.Sprintf("task-%d", t.counter), 5, rand.IntN(4), t.wg)
}
