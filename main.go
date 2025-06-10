package main

import (
	"flag"
	"fmt"
	"practice/assignments/logic"
)

func main() {
	add := flag.String("add", "", "Add a new todo item")
	deleteID := flag.Int("delete", -1, "Delete a todo item by ID")
	update := flag.Int("update", -1, "Update the description or status of a todo item by ID")
	newDesc := flag.String("desc", "", "New description for update")
	list := flag.Bool("list", false, "List all items")
	status := flag.String("status", "", "New status for the item (not started, started, completed)")
	flag.Parse()

	todos := logic.LoadTodos()

	if *add != "" {
		todos = logic.AddTodo(todos, *add)
		fmt.Println("New item added.")
		fmt.Println("--------------------")
	}

	if *deleteID != -1 {
		found := false
		for i := 0; i < len(todos); i++ {
			if todos[i].ID == *deleteID {
				todos = append(todos[:i], todos[i+1:]...)
				fmt.Println("Item deleted.")
				fmt.Println("--------------------")
				found = true
				break
			}
		}
		if !found {
			fmt.Println("--------------------")
			fmt.Println("Item not found.")
		}
	}

	if *list || *add != "" {
		fmt.Println("To-Do List:")
		fmt.Println("--------------------")
		for _, t := range todos {
			fmt.Printf("%d: %s [%s]\n", t.ID, t.Description, t.Status)
		}
	}

	if *update != -1 && *status != "" {
		for i := range todos {
			if todos[i].ID == *update {
				todos[i].Status = *status
				fmt.Println("Status Updated.")
				fmt.Println("--------------------")
			}
		}
	}

	if *update != -1 && *newDesc != "" {
		for i := range todos {
			if todos[i].ID == *update {
				todos[i].Description = *newDesc
				fmt.Println("Description updated.")
				fmt.Println("--------------------")
			}
		}
	}

	logic.SaveTodos(todos)
}
