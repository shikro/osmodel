package processor_test

import (
	"context"
	"testing"
	"time"

	"github.com/shikro/osmodel/processor"
	"github.com/shikro/osmodel/scheduler"
	"github.com/shikro/osmodel/task"
	"github.com/stretchr/testify/assert"
)

func TestProcessor(t *testing.T) {
	t.Run("Should execute task and notify scheduler when task done", func(t *testing.T) {
		ctx := context.Background()
		p := processor.New()
		s := scheduler.MockScheduler{}
		p.SetScheduler(&s)

		task := &task.MockTask{FinishTask: true}

		p.Start(ctx)
		p.ExecuteTask(task)
		time.Sleep(100 * time.Millisecond)

		assert.True(t, task.DoWorkCalled)
		assert.True(t, s.TaskDoneCalled)
	})

	t.Run("Should notify scheduler when task replaced", func(t *testing.T) {
		ctx := context.Background()
		p := processor.New()
		s := scheduler.MockScheduler{}
		p.SetScheduler(&s)

		t1 := &task.MockTask{}
		t2 := &task.MockTask{}

		p.Start(ctx)
		p.ExecuteTask(t1)
		p.ExecuteTask(t2)
		time.Sleep(100 * time.Millisecond)

		assert.True(t, s.RetakeTaskCalled)
	})

	t.Run("Should notify scheduler when task waiting", func(t *testing.T) {
		ctx := context.Background()
		p := processor.New()
		s := scheduler.MockScheduler{}
		p.SetScheduler(&s)

		task := &task.MockTask{NeedToWait: true}

		p.Start(ctx)
		p.ExecuteTask(task)
		time.Sleep(100 * time.Millisecond)

		assert.True(t, task.WaitCalled)
		assert.True(t, s.TaskWaitingCalled)
	})
}
