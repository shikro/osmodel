package task

type MockTask struct {
	Id           string
	Prior        int
	DoWorkCalled bool
	WaitCalled   bool
	FinishTask   bool
	NeedToWait   bool
}

func (t MockTask) ID() string {
	return t.Id
}

func (t *MockTask) DoWork() {
	t.DoWorkCalled = true
}

func (t MockTask) Finished() bool {
	return t.FinishTask
}

func (t MockTask) Priority() int {
	return t.Prior
}

func (t *MockTask) Wait() bool {
	t.WaitCalled = true
	return t.NeedToWait
}
