package data

// StoreInterface defines the operations that can be performed on a task store.
type StoreInterface interface {
	// AddTask adds a task to the store and assigns it a unique ID.
	// Returns an error if the operation fails.
	AddTask(t *Task) error

	// ListAllTasks returns all tasks in the store.
	ListAllTasks() []*Task

	// GetTaskByID retrieves a task by its ID.
	// Returns nil if no task with the given ID exists.
	GetTaskByID(id int) *Task

	// UpdateTask updates an existing task.
	// Returns an error if the task doesn't exist or the operation fails.
	UpdateTask(t *Task) error

	// DeleteTask removes a task from the store.
	// Returns an error if the task doesn't exist or the operation fails.
	DeleteTask(id int) error
}
