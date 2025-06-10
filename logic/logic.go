package logic

import (
	model "practice/assignments/model"
	"practice/assignments/storage"
)

func LoadTodos() []model.TodoItem {
	return storage.LoadTodos()
}

func SaveTodos(todos []model.TodoItem) {
	storage.SaveTodos(todos)
}

func AddTodo(todos []model.TodoItem, desc string) []model.TodoItem {
	newID := getNextID(todos)
	todos = append(todos, model.TodoItem{
		ID:          newID,
		Description: desc,
		Status:      "not started",
	})
	return todos
}

func getNextID(todos []model.TodoItem) int {
	maxID := 0
	for _, item := range todos {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}
