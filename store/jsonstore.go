package store

import (
	"encoding/json"
	"fmt"
	"github.com/kevin7254/task/model"
	"os"
	"path/filepath"
	"sync"
)

// JsonStore implements the model.TaskRepository for persisting tasks to a JSON file.
type JsonStore struct {
	filename string
	tasks    map[int]*model.Task
	mu       sync.RWMutex
	nextID   int
}

// NewJsonStore creates a new JsonStore instance that persists tasks to the specified file.
// It creates the directory if it doesn't exist and loads any existing tasks.
func NewJsonStore(filename string) (*JsonStore, error) {
	if osErr := os.MkdirAll(filepath.Dir(filename), 0755); osErr != nil {
		return nil, fmt.Errorf("failed to create directory: %w", osErr)
	}

	jsonStore := &JsonStore{
		filename: filename,
		tasks:    make(map[int]*model.Task),
		nextID:   1,
	}

	if jsonErr := jsonStore.load(); jsonErr != nil {
		if !os.IsNotExist(jsonErr) {
			return nil, fmt.Errorf("failed to load tasks: %w", jsonErr)
		}
		// File doesn't exist yet, which is fine for a new store
	} else {
		jsonStore.updateNextID()
	}

	return jsonStore, nil
}

// updateNextID sets the nextID to be one more than the highest ID in the tasks map.
func (s *JsonStore) updateNextID() {
	highestID := 0
	for id := range s.tasks {
		if id > highestID {
			highestID = id
		}
	}
	s.nextID = highestID + 1
}

// save persists the tasks to the store file.
func (s *JsonStore) save() error {
	s.mu.RLock()
	bytes, marshErr := json.MarshalIndent(s.tasks, "", "  ")
	s.mu.RUnlock()

	if marshErr != nil {
		return fmt.Errorf("failed to marshal tasks: %w", marshErr)
	}

	if osErr := os.WriteFile(s.filename, bytes, 0644); osErr != nil {
		return fmt.Errorf("failed to write tasks to file: %w", osErr)
	}

	return nil
}

// load reads tasks from the store file.
func (s *JsonStore) load() error {
	bytes, osErr := os.ReadFile(s.filename)
	if osErr != nil {
		return osErr
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if unMarshalErr := json.Unmarshal(bytes, &s.tasks); unMarshalErr != nil {
		return fmt.Errorf("failed to unmarshal tasks: %w", unMarshalErr)
	}

	return nil
}

// ListAllTasks returns all tasks in the store.
func (s *JsonStore) ListAllTasks() []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

// GetTaskByID retrieves a task by its ID.
// Returns nil if no task with the given ID exists.
func (s *JsonStore) GetTaskByID(id int) *model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.tasks[id]
}

// AddTask adds a task to the store and assigns it a unique ID.
// Returns an error if the operation fails.
func (s *JsonStore) AddTask(t *model.Task) error {
	s.mu.Lock()
	t.ID = s.nextID
	s.nextID++
	s.tasks[t.ID] = t
	s.mu.Unlock()

	return s.save()
}

// UpdateTask updates an existing task.
// Returns an error if the task doesn't exist or the operation fails.
func (s *JsonStore) UpdateTask(t *model.Task) error {
	s.mu.Lock()

	if _, exists := s.tasks[t.ID]; !exists {
		return fmt.Errorf("task with ID %d does not exist", t.ID)
	}

	s.tasks[t.ID] = t
	s.mu.Unlock()
	return s.save()
}

// DeleteTask removes a task from the store.
// Returns an error if the operation fails.
func (s *JsonStore) DeleteTask(id int) error {
	s.mu.Lock()
	delete(s.tasks, id)
	s.mu.Unlock()

	return s.save()
}
