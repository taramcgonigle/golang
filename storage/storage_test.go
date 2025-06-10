package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"practice/assignments/model"
)

// Helper to create a temporary test file with content
func createTempTodoFile(t *testing.T, todos []model.TodoItem) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "todo-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	data, _ := json.MarshalIndent(todos, "", "  ")
	_, err = tmpFile.Write(data)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestLoadTodosFromValidFile(t *testing.T) {
	expected := []model.TodoItem{
		{ID: 1, Description: "Test task", Status: "started"},
	}

	// Create temp file with JSON content
	tempFile := createTempTodoFile(t, expected)
	defer os.Remove(tempFile)

	// Override the const fileName for testing
	original := fileName
	fileName = tempFile
	defer func() { fileName = original }()

	result := LoadTodos()

	if len(result) != 1 || result[0].Description != expected[0].Description {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestLoadTodosFileNotFound(t *testing.T) {
	original := fileName
	fileName = "nonexistent.json"
	defer func() { fileName = original }()

	result := LoadTodos()
	if len(result) != 0 {
		t.Errorf("Expected empty slice on missing file, got %+v", result)
	}
}

func TestSaveTodosCreatesValidFile(t *testing.T) {
	expected := []model.TodoItem{
		{ID: 1, Description: "Saved task", Status: "completed"},
	}

	// Create temp file path (not open)
	tempFile := filepath.Join(os.TempDir(), "todo-test-save.json")
	defer os.Remove(tempFile)

	// Redirect the fileName for test
	original := fileName
	fileName = tempFile
	defer func() { fileName = original }()

	SaveTodos(expected)

	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Could not read saved file: %v", err)
	}

	var decoded []model.TodoItem
	err = json.Unmarshal(content, &decoded)
	if err != nil || decoded[0].Description != expected[0].Description {
		t.Errorf("Saved content mismatch: got %+v, expected %+v", decoded, expected)
	}
}
