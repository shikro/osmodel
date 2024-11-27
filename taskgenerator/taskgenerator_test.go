package taskgenerator_test

import (
	"sync"
	"testing"

	"github.com/shikro/osmodel/taskgenerator"
	"github.com/stretchr/testify/assert"
)

func TestTaskGenerator(t *testing.T) {
	t.Run("Should create task", func(t *testing.T) {
		var wg sync.WaitGroup
		gen := taskgenerator.New(&wg)
		task := gen.Next()
		assert.NotPanics(t, wg.Done)
		assert.NotEmpty(t, task.ID())
		assert.True(t, task.Priority() > -1 && task.Priority() < 4)
		assert.False(t, task.Finished())
	})
}
