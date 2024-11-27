package scheduler_test

import (
	"context"
	"sync"
	"testing"

	"github.com/shikro/osmodel/processor"
	"github.com/shikro/osmodel/scheduler"
	"github.com/shikro/osmodel/task"
	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	t.Run("Should execute tasks with high prior first", func(t *testing.T) {
		ctx := context.Background()
		s := scheduler.New(func() chan struct{} { return make(chan struct{}) })
		var wg sync.WaitGroup
		p := processor.MockProcessor{Wg: &wg}
		s.SetProcessor(&p)

		wg.Add(4)

		t1 := &task.MockTask{Id: "0", Prior: 0}
		t2 := &task.MockTask{Id: "1", Prior: 1}
		t3 := &task.MockTask{Id: "2", Prior: 2}
		t4 := &task.MockTask{Id: "3", Prior: 3}

		s.Start(ctx)
		s.ScheldueTask(t1)
		s.ScheldueTask(t3)
		s.ScheldueTask(t4)
		s.ScheldueTask(t2)
		s.TaskDone()
		s.TaskDone()
		s.TaskDone()
		s.TaskDone()

		wg.Wait()

		assert.Equal(t, t1.Id, p.Tasks[0].ID())
		assert.Equal(t, t2.Id, p.Tasks[1].ID())
		assert.Equal(t, t3.Id, p.Tasks[2].ID())
		assert.Equal(t, t4.Id, p.Tasks[3].ID())
	})

	t.Run("Should give task after retake first", func(t *testing.T) {
		ctx := context.Background()
		s := scheduler.New(func() chan struct{} { return make(chan struct{}) })
		var wg sync.WaitGroup
		p := processor.MockProcessor{Wg: &wg}
		s.SetProcessor(&p)

		wg.Add(4)

		t1 := &task.MockTask{Id: "0", Prior: 0}
		t2 := &task.MockTask{Id: "1", Prior: 1}
		t3 := &task.MockTask{Id: "2", Prior: 1}

		s.Start(ctx)
		s.ScheldueTask(t2)
		s.ScheldueTask(t3)
		s.ScheldueTask(t1)

		s.RetakeTask(t2)
		s.TaskDone()
		s.TaskDone()
		s.TaskDone()

		wg.Wait()

		assert.Equal(t, t2.Id, p.Tasks[0].ID())
		assert.Equal(t, t1.Id, p.Tasks[1].ID())
		assert.Equal(t, t2.Id, p.Tasks[2].ID())
		assert.Equal(t, t3.Id, p.Tasks[3].ID())
	})

	t.Run("Should awake task after event", func(t *testing.T) {
		ctx := context.Background()
		event := make(chan struct{})
		s := scheduler.New(func() chan struct{} { return event })
		var wg sync.WaitGroup
		p := processor.MockProcessor{Wg: &wg}
		s.SetProcessor(&p)

		wg.Add(2)

		task := &task.MockTask{Id: "0", Prior: 0}

		s.Start(ctx)
		s.ScheldueTask(task)

		s.TaskWaiting(task)
		event <- struct{}{}
		s.TaskDone()

		wg.Wait()

		assert.Equal(t, 2, len(p.Tasks))
	})

}
