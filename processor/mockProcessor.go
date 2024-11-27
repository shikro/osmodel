package processor

import (
	"sync"

	"github.com/shikro/osmodel/task"
)

type MockProcessor struct {
	Tasks []task.Task
	Wg    *sync.WaitGroup
}

func (p *MockProcessor) ExecuteTask(task task.Task) {
	p.Tasks = append(p.Tasks, task)
	p.Wg.Done()
}
