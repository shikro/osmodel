package scheduler

import "github.com/shikro/osmodel/task"

type MockScheduler struct {
	TaskDoneCalled    bool
	RetakeTaskCalled  bool
	TaskWaitingCalled bool
}

func (s *MockScheduler) RetakeTask(task task.Task) {
	s.RetakeTaskCalled = true
}

func (s *MockScheduler) TaskDone() {
	s.TaskDoneCalled = true

}

func (s *MockScheduler) TaskWaiting(task task.Task) {
	s.TaskWaitingCalled = true
}
