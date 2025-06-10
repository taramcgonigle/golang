package logic

import (
	"testing"

	"practice/assignments/model"
)

func TestAddTodo(t *testing.T) {
	todos := []model.TodoItem{}

	todos = AddTodo(todos, "Test adding")

	if len(todos) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(todos))
	}

	if todos[0].ID != 1 {
		t.Errorf("Expected ID 1, got %d", todos[0].ID)
	}

	if todos[0].Description != "Test adding" {
		t.Errorf("Expected description 'Test adding', got '%s'", todos[0].Description)
	}

	if todos[0].Status != "not started" {
		t.Errorf("Expected status 'not started', got '%s'", todos[0].Status)
	}
}

func TestAddMultipleTodos(t *testing.T) {
	todos := []model.TodoItem{
		{ID: 1, Description: "First", Status: "not started"},
		{ID: 2, Description: "Second", Status: "completed"},
	}

	todos = AddTodo(todos, "Third")

	if len(todos) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(todos))
	}

	if todos[2].ID != 3 {
		t.Errorf("Expected new item ID 3, got %d", todos[2].ID)
	}
}

func TestGetNextID(t *testing.T) {
	cases := []struct {
		name     string
		todos    []model.TodoItem
		expected int
	}{
		{"Empty list", []model.TodoItem{}, 1},
		{"Single item", []model.TodoItem{{ID: 5}}, 6},
		{"Unordered IDs", []model.TodoItem{{ID: 3}, {ID: 7}, {ID: 2}}, 8},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getNextID(c.todos)
			if got != c.expected {
				t.Errorf("Expected next ID %d, got %d", c.expected, got)
			}
		})
	}
}
