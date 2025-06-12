package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"practice/assignments/model"
)

// Helper to create temporary test file with content
func createTempTodoFile(t *testing.T, todos []model.TodoItem) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "todo-test-*.json") // Creates temporary file in the OS's temp directory
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err) // Fatal if file creation fail
	}

	//add todos to Json file
	data, _ := json.MarshalIndent(todos, "", "  ")
	_, err = tmpFile.Write(data)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

// test loading todos from existing file
func TestLoadTodosFromValidFile(t *testing.T) {
	expected := []model.TodoItem{
		{ID: 1, Description: "Test task", Status: "started"},
	}

	tempFile := createTempTodoFile(t, expected) // Create temp file with JSON content
	defer os.Remove(tempFile)

	// Override the const fileName for testing
	original := fileName
	fileName = tempFile
	defer func() { fileName = original }() // Restore original fileName after test

	result := LoadTodos() // Call the function under test

	if len(result) != 1 || result[0].Description != expected[0].Description {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

// testing loading todos from an empty file
func TestLoadTodosFileNotFound(t *testing.T) {
	original := fileName // Point to a non-existent file
	fileName = "nonexistent.json"
	defer func() { fileName = original }()

	// This should return an empty slice instead of crashing
	result := LoadTodos()
	if len(result) != 0 {
		t.Errorf("Expected empty slice on missing file, got %+v", result)
	}
}

// Test that SaveTodos writes valid JSON to disk
func TestSaveTodosCreatesValidFile(t *testing.T) {
	expected := []model.TodoItem{
		{ID: 1, Description: "Saved task", Status: "completed"},
	}

	tempFile := filepath.Join(os.TempDir(), "todo-test-save.json") // Create path in temp directory to write to
	defer os.Remove(tempFile)

	// Override the storage file path
	original := fileName
	fileName = tempFile
	defer func() { fileName = original }()

	SaveTodos(expected) // Call SaveTodos, which should write to the temp file

	content, err := os.ReadFile(tempFile) // Read the file back and parse it
	if err != nil {
		t.Fatalf("Could not read saved file: %v", err)
	}

	var decoded []model.TodoItem
	err = json.Unmarshal(content, &decoded)
	if err != nil || decoded[0].Description != expected[0].Description {
		t.Errorf("Saved content mismatch: got %+v, expected %+v", decoded, expected)
	}
}
