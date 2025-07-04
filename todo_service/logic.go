package logic

import (
	model "practice/assignments/model"
	"practice/assignments/storage"
)

// LoadTodos loads the current list of to-do items from storage
func LoadTodos() []model.TodoItem {
	return storage.LoadTodos()
}

// SaveTodos writes the given list of to-do items to storage
func SaveTodos(todos []model.TodoItem) {
	storage.SaveTodos(todos)
}

// AddTodo appends a new to-do item with a generated unique ID and a default "not started" status to the given slice of todos
func AddTodo(todos []model.TodoItem, desc string) []model.TodoItem {
	newID := getNextID(todos)
	todos = append(todos, model.TodoItem{
		ID:          newID,
		Description: desc,
		Status:      "not started",
	})
	return todos // Returns a new slice with the added item.
}

// getNextID scans the current list of to-dos to find the highest ID, then returns the next available ID
func getNextID(todos []model.TodoItem) int {
	maxID := 0
	for _, item := range todos {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}
