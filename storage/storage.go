package storage

import (
	"encoding/json"
	"os"
	"practice/assignments/model"
)

var fileName = "todo.json"

func SaveTodos(todos []model.TodoItem) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func LoadTodos() []model.TodoItem {
	var todos []model.TodoItem
	file, err := os.ReadFile(fileName)
	if err != nil {
		return todos
	}
	_ = json.Unmarshal(file, &todos)
	return todos
}
