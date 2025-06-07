package model

import (
	"fmt"
	"time"
)

type Priority int

const (
	Low Priority = iota + 1
	Medium
	High
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Project     string    `json:"project"`
	Priority    Priority  `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	TimeSpent   int64     `json:"time_spent"`
}

func NewTask(title string, description string, project string, priority Priority, dueDate time.Time) *Task {
	return &Task{
		ID:          0,
		Title:       title,
		Description: description,
		Project:     project,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
	}
}

func (t *Task) String() string {
	status := "⏳"
	if !t.CompletedAt.IsZero() {
		status = "✅"
	} else if t.IsOverdue() {
		status = "⚠️"
	}

	return fmt.Sprintf("%s %s [%s] Due: %s (%d min spent)",
		status,
		t.Title,
		t.Project,
		t.DueDate.Format("2006-01-02"),
		t.TimeSpent,
	)
}

func (t *Task) IsOverdue() bool {
	return t.DueDate.Before(time.Now())
}

func (t *Task) Complete() {
	t.CompletedAt = time.Now()
}

func (t *Task) AddTimeSpent(minutes int64) {
	t.TimeSpent += minutes
}
