package task_test

import (
	"sync"
	"testing"

	"github.com/shikro/osmodel/task"
	"github.com/stretchr/testify/assert"
)

func TestTasks(t *testing.T) {
	t.Run("Should return ID of task", func(t *testing.T) {
		task := task.NewSimpleTask("id", 1, 1, nil)
		assert.Equal(t, "id", task.ID())
	})

	t.Run("Should return priority of task", func(t *testing.T) {
		task := task.NewSimpleTask("id", 1, 1, nil)
		assert.Equal(t, 1, task.Priority())
	})

	t.Run("Should do work and finish task", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		task := task.NewSimpleTask("id", 1, 1, &wg)
		assert.False(t, task.Finished())
		task.DoWork()
		assert.True(t, task.Finished())
	})

	t.Run("Should start waiting single time", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		task := task.NewTaskWithWait("id", 2, 1, &wg)
		assert.False(t, task.Wait())
		task.DoWork()
		assert.True(t, task.Wait())
		assert.False(t, task.Wait())
	})
}
