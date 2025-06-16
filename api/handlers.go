package api

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"practice/assignments/model"
	"practice/assignments/storage"
	"practice/assignments/trace"
	"strconv"
)

// get handler serves the homepage which displays the list of todos
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	traceID := trace.GetTraceID(r.Context()) //log the request trace ID
	slog.Info("Serving homepage", "traceID", traceID)

	todos := storage.LoadTodos() //load all current todos from json file

	//redner the homepage template + passing the todos
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	err := tmpl.Execute(w, todos)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// get /create-form
func CreateFormHandler(w http.ResponseWriter, r *http.Request) {
	// Render the create page template
	tmpl := template.Must(template.ParseFiles("templates/create.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering form", http.StatusInternalServerError)
	}
}

// post /create
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	traceID := trace.GetTraceID(r.Context())

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	//parse the form input from the create page
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	todos := storage.LoadTodos() //Load current list of todos

	//generate a new ID based on last item's ID
	newID := 0
	if len(todos) > 0 {
		newID = todos[len(todos)-1].ID + 1
	}

	//build a new TodoItem from form values
	todo := model.TodoItem{
		ID:          newID,
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
	}

	//Add to list and save to file
	todos = append(todos, todo)
	err := storage.SaveTodos(todos)
	if err != nil {
		http.Error(w, "Failed to save todo", http.StatusInternalServerError)
		return
	}

	slog.Info("Todo created", "traceID", traceID, "ID", todo.ID, "description", todo.Description)

	http.Redirect(w, r, "/", http.StatusSeeOther) //go back to homepage
}

// Get
func GetHandler(w http.ResponseWriter, r *http.Request) {
	traceID := trace.GetTraceID(r.Context())
	slog.Info("Returning JSON todo list", "traceID", traceID)

	todos := storage.LoadTodos()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// get or post /edit-form
func EditFormHandler(w http.ResponseWriter, r *http.Request) {
	var idStr string

	// Handle both GET and POST submissions of selected ID
	if r.Method == http.MethodPost {
		idStr = r.FormValue("id")
	} else {
		idStr = r.URL.Query().Get("id")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todos := storage.LoadTodos()

	// Find the selected todo item
	var selected model.TodoItem
	for _, t := range todos {
		if t.ID == id {
			selected = t
			break
		}
	}

	// Render the edit page with the selected item
	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	err = tmpl.Execute(w, selected)
	if err != nil {
		http.Error(w, "Error rendering edit page", http.StatusInternalServerError)
	}
}

// post /update
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	desc := r.FormValue("description")
	status := r.FormValue("status")

	traceID := trace.GetTraceID(r.Context())
	slog.Info("Update requested", "traceID", traceID, "id", idStr, "description", desc, "status", status)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todos := storage.LoadTodos()
	updated := false

	// Replace the matching item in the list
	for i, t := range todos {
		if t.ID == id {
			todos[i].Description = desc
			todos[i].Status = status
			updated = true
			break
		}
	}

	if !updated {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if err := storage.SaveTodos(todos); err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	slog.Info("Todo updated", "traceID", traceID, "id", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// post /delete
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	traceID := trace.GetTraceID(r.Context())
	idStr := r.FormValue("id") // from form body (not URL query)

	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todos := storage.LoadTodos()
	index := -1
	for i := range todos {
		if todos[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Remove the item from the slice
	todos = append(todos[:index], todos[index+1:]...)
	_ = storage.SaveTodos(todos)

	slog.Info("Todo deleted", "traceID", traceID, "id", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
