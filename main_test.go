package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func Test_inArray(t *testing.T) {
	t.Run("is serach term in array", func(t *testing.T) {
		list := []string{"hello", "world", "!"}
		term := "world"
		want := true

		got := inArray(term, list...)

		if got != want {
			t.Errorf("inArray() = %v, want %v", got, want)
		}
	})

	t.Run("search for item in array return false", func(t *testing.T) {
		list := []string{"hello", "world", "!"}
		term := "john"
		want := false

		got := inArray(term, list...)
		if got != want {
			t.Errorf("inArray() returns %v, wants %v", got, want)
		}
	})
}

func Test_parseEqualTo(t *testing.T) {
	t.Run("remove all = in the list", func(t *testing.T) {
		list := []string{"a=b", "c", "=", "d=e=f", "g", "  q=v  ", "volvo", " = ", "d=e=a=f_jam="}
		want := []string{"a", "b", "c", "d", "e", "f", "g", "q", "v", "volvo", "d", "e", "a", "f_jam"}

		got := parseEqualTo(list)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("parseEqualTo() = %v, want %v", got, want)
		}
	})

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input    []string
		expected []string
	}{
		{
			name:     "Basic key-value pair",
			input:    []string{"a=b"},
			expected: []string{"a", "b"},
		},
		{
			name:     "Multiple key-value pairs",
			input:    []string{"a=b", "c=d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "Mixed list with spaces",
			input:    []string{"  x=y  ", "  z  ", "  =  "},
			expected: []string{"x", "y", "z"},
		},
		{
			name:     "String with multiple equals signs",
			input:    []string{"key=value=more"},
			expected: []string{"key", "value", "more"},
		},
		{
			name:     "Empty list",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "List with only equals signs",
			input:    []string{"=", " = "},
			expected: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseEqualTo(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parseEqualTo() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_checkArguementType(t *testing.T) {
	// Setup: Create a temporary directory and a temporary file for testing
	tempDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile, err := os.CreateTemp("", "test-file-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close() // Close the file to allow for cleanup
	defer os.Remove(tempFilePath)

	// Define test cases in a table-driven format
	tests := []struct {
		name      string
		arguement interface{}
		argType   string
		wantErr   bool
		errString string
	}{
		// --- "directory" Test Cases ---
		{
			name:      "Valid directory path",
			arguement: tempDir,
			argType:   "directory",
			wantErr:   false,
		},
		{
			name:      "Directory path does not exist",
			arguement: filepath.Join(tempDir, "non-existent-dir"),
			argType:   "directory",
			wantErr:   true,
			errString: "Path does not exist",
		},
		{
			name:      "Path is a file, not a directory",
			arguement: tempFilePath,
			argType:   "directory",
			wantErr:   true,
			errString: "path is not a directory",
		},
		{
			name:      "Incorrect type for directory (int)",
			arguement: 123,
			argType:   "directory",
			wantErr:   true,
			errString: "expected a string for directory path",
		},
		{
			name:      "Incorrect type for directory (bool)",
			arguement: true,
			argType:   "directory",
			wantErr:   true,
			errString: "expected a string for directory path",
		},

		// --- "boolean" Test Cases ---
		{
			name:      "Valid boolean (true)",
			arguement: true,
			argType:   "boolean",
			wantErr:   false,
		},
		{
			name:      "Valid boolean (false)",
			arguement: false,
			argType:   "boolean",
			wantErr:   false,
		},
		{
			name:      "Incorrect type for boolean (string)",
			arguement: "true",
			argType:   "boolean",
			wantErr:   true,
			errString: "expected a boolean for: true",
		},
		{
			name:      "Incorrect type for boolean (int)",
			arguement: 1,
			argType:   "boolean",
			wantErr:   true,
			errString: "expected a boolean for: 1",
		},

		// --- "default" Test Cases ---
		{
			name:      "Unsupported type string",
			arguement: "some-value",
			argType:   "unsupported",
			wantErr:   true,
			errString: "unspecified type received for: some-value",
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkArgumentType(tt.arguement, tt.argType)

			// Check if the error presence matches the expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("checkArguementType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If an error is expected, check if the error message is correct
			if tt.wantErr && tt.errString != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errString) {
					t.Errorf("checkArguementType() error message got = %q, want contains %q", err, tt.errString)
				}
			}
		})
	}
}

func Test_parseFlagsAndArgument(t *testing.T) {
	// Setup: Create a temporary directory and a temporary file for testing
	tempDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile, err := os.CreateTemp("", "test-file-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close() // Close the file to allow for cleanup
	defer os.Remove(tempFilePath)

	destinationTempDir, err := os.MkdirTemp("", "test-dir-destination-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(destinationTempDir)

	absPathDest, err := filepath.Abs(destinationTempDir)
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}

	absPathTemp, err := filepath.Abs(tempDir)
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}

	tests := []struct {
		name string
		args []string
		want map[string]any
	}{
		{
			name: "Success with mixed flags",
			args: []string{"--dry-run", "--source", tempDir, "--destination", destinationTempDir},
			want: map[string]any{"dry-run": true, "dotfiles": false, "source": absPathTemp, "destination": absPathDest},
		},
		{
			name: "Success with target and mixed flags",
			args: []string{"--dry-run", "--target", tempDir, "--dotfiles", "--destination", destinationTempDir},
			want: map[string]any{"dry-run": true, "dotfiles": true, "target": absPathTemp, "destination": absPathDest},
		},
		{
			name: "Only boolean flags",
			args: []string{"--dry-run", "--dotfiles"},
			want: map[string]any{"dry-run": true, "dotfiles": true},
		},
		{
			name: "No flags",
			args: []string{"arg1", "arg2"},
			want: map[string]any{"dry-run": false, "dotfiles": false},
		},
		{
			name: "Empty argument list",
			args: []string{},
			want: map[string]any{"dry-run": false, "dotfiles": false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect log.Fatalf to a custom function to prevent the test from exiting.
			// This allows us to test the happy path correctly.
			var got map[string]any

			// This defer/recover block is not perfect for `log.Fatalf` which calls os.Exit,
			// but it demonstrates the intent to test for panics in a test environment.
			defer func() {
				if r := recover(); r != nil {
					// We panic only for fatal errors, which are not expected in these tests.
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			got = parseFlagsAndArgument(tt.args)

			// Use reflect.DeepEqual for a robust map comparison.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFlagsAndArgument() got = %v, want %v", got, tt.want)
			}
		})
	}
}
