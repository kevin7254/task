package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"task/data"
)

// Storage implements the data.StoreInterface for persisting tasks to a JSON file.
type Storage struct {
	filename string
	tasks    map[int]*data.Task
	mu       sync.RWMutex
	nextID   int
}

// NewStorage creates a new Storage instance that persists tasks to the specified file.
// It creates the directory if it doesn't exist and loads any existing tasks.
func NewStorage(filename string) (*Storage, error) {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	s := &Storage{
		filename: filename,
		tasks:    make(map[int]*data.Task),
		nextID:   1,
	}

	if err := s.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load tasks: %w", err)
		}
		// File doesn't exist yet, which is fine for a new storage
	} else {
		// Determine the next available ID
		s.updateNextID()
	}

	return s, nil
}

// updateNextID sets the nextID to be one more than the highest ID in the tasks map.
func (s *Storage) updateNextID() {
	highestID := 0
	for id := range s.tasks {
		if id > highestID {
			highestID = id
		}
	}
	s.nextID = highestID + 1
}

// save persists the tasks to the storage file.
func (s *Storage) save() error {
	s.mu.RLock()
	bytes, err := json.MarshalIndent(s.tasks, "", "  ")
	s.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := os.WriteFile(s.filename, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write tasks to file: %w", err)
	}

	return nil
}

// load reads tasks from the storage file.
func (s *Storage) load() error {
	bytes, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := json.Unmarshal(bytes, &s.tasks); err != nil {
		return fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return nil
}

// ListAllTasks returns all tasks in the store.
func (s *Storage) ListAllTasks() []*data.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*data.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

// GetTaskByID retrieves a task by its ID.
// Returns nil if no task with the given ID exists.
func (s *Storage) GetTaskByID(id int) *data.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.tasks[id]
}

// AddTask adds a task to the store and assigns it a unique ID.
// Returns an error if the operation fails.
func (s *Storage) AddTask(t *data.Task) error {
	s.mu.Lock()
	t.ID = s.nextID
	s.nextID++
	s.tasks[t.ID] = t
	s.mu.Unlock()

	return s.save()
}

// UpdateTask updates an existing task.
// Returns an error if the task doesn't exist or the operation fails.
func (s *Storage) UpdateTask(t *data.Task) error {
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
func (s *Storage) DeleteTask(id int) error {
	s.mu.Lock()
	delete(s.tasks, id)
	s.mu.Unlock()

	return s.save()
}
