package storage

import (
	"encoding/json"
	"log/slog"
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
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal todos", "error", err)
		return
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		slog.Error("Failed to write file", "error", err)
		return
	}
	slog.Info("Todos saved", "count", len(todos))
}
