package storage

import (
	"encoding/json"
	"os"
	"practice/assignments/model"
)

var fileName = "todo.json" //name of the file where todos are stored

// serializes a slice of TodoItem Structs to JSON and saves it to a file
func SaveTodos(todos []model.TodoItem) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil { //If serialization fails, the error is returned to the caller.
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

// LoadTodos attempts to read the todo.json file and parse its content into a slice of TodoItem structs
func LoadTodos() []model.TodoItem {
	var todos []model.TodoItem
	file, err := os.ReadFile(fileName)
	if err != nil { // If the file does not exist or fails to parse, an empty list is returned
		return todos
	}
	_ = json.Unmarshal(file, &todos) // Ignore unmarshal errors for simplicity
	return todos
}

// DeleteByIndex removes the item at the given index from the slice and returns the updated slice
func DeleteByIndex(todos []model.TodoItem, index int) ([]model.TodoItem, bool) {
	if index < 0 || index >= len(todos) {
		return todos, false
	}
	//Rebuilds the slice without the deleted item included
	todos = append(todos[:index], todos[index+1:]...)
	return todos, true
}
