package storage

import (
	"encoding/json"
	"os"
	model "practice/assignments/model"
)

var fileName = "todo.json"

func LoadTodos() []model.TodoItem {
	var todos []model.TodoItem
	file, err := os.ReadFile(fileName)
	if err != nil {
		return todos
	}
	_ = json.Unmarshal(file, &todos)
	return todos
}

func SaveTodos(todos []model.TodoItem) {
	data, _ := json.MarshalIndent(todos, "", "  ")
	_ = os.WriteFile(fileName, data, 0644)
}
