package logic

import (
	"testing"

	"practice/assignments/model"
)

// TestAddTodo verifies that a new todo item is correctly added to an empty list. Checks item has correct ID, description, and status.
func TestAddTodo(t *testing.T) {
	todos := []model.TodoItem{}

	todos = AddTodo(todos, "Test adding")

	if len(todos) != 1 { // Verify item 1 was added correctly
		t.Fatalf("Expected 1 item, got %d", len(todos))
	}

	if todos[0].ID != 1 { // Check that ID was initialized as 1
		t.Errorf("Expected ID 1, got %d", todos[0].ID)
	}

	if todos[0].Description != "Test adding" { // Validate correct description
		t.Errorf("Expected description 'Test adding', got '%s'", todos[0].Description)
	}

	if todos[0].Status != "not started" { // Ensure default status is set properly
		t.Errorf("Expected status 'not started', got '%s'", todos[0].Status)
	}
}

// TestUpdateTodo verifies that an existing todo item can be updated with a new description and status.
func TestAddMultipleTodos(t *testing.T) {
	todos := []model.TodoItem{
		{ID: 1, Description: "First", Status: "not started"},
		{ID: 2, Description: "Second", Status: "completed"},
	}

	todos = AddTodo(todos, "Third")

	if len(todos) != 3 { // Check that the new item was added correctly
		t.Fatalf("Expected 3 items, got %d", len(todos))
	}

	if todos[2].ID != 3 { // Verify that the new item has the correct ID automatically incremented
		t.Errorf("Expected new item ID 3, got %d", todos[2].ID)
	}
}

// tests to validate getNextID across multiple scenarios
func TestGetNextID(t *testing.T) {
	cases := []struct {
		name     string
		todos    []model.TodoItem
		expected int
	}{
		{"Empty list", []model.TodoItem{}, 1},                             // No items yet, start from 1
		{"Single item", []model.TodoItem{{ID: 5}}, 6},                     // Max ID is 5, expect 6
		{"Unordered IDs", []model.TodoItem{{ID: 3}, {ID: 7}, {ID: 2}}, 8}, // Max ID is 7, expect 8
	}

	// Run each case using subtests for clarity
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getNextID(c.todos)
			if got != c.expected {
				t.Errorf("Expected next ID %d, got %d", c.expected, got)
			}
		})
	}
}
