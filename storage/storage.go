package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"task/data"
)

type Storage struct {
	filename string
	tasks    map[int]*data.Task
	mu       sync.RWMutex
}

func NewStorage(filename string) (*Storage, error) {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	s := &Storage{
		filename: filename,
		tasks:    make(map[int]*data.Task),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load tasks %w", err)
	}
	return s, nil
}

func (s *Storage) save() error {
	s.mu.RLock()
	bytes, err := json.MarshalIndent(s.tasks, "", "  ")
	s.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal tasks %w", err)
	}
	return os.WriteFile(s.filename, bytes, 0644)
}

func (s *Storage) load() error {
	bytes, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	return json.Unmarshal(bytes, &s.tasks)
}

func (s *Storage) ListAllTasks() []*data.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*data.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *Storage) AddTask(t *data.Task) error {
	s.mu.Lock()
	s.tasks[t.ID] = t
	s.mu.Unlock()

	return s.save()
}

func (s *Storage) DeleteTask(id int) error {
	s.mu.Lock()
	delete(s.tasks, id)
	s.mu.Unlock()

	return s.save()
}
