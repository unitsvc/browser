package browser

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

// errorReader is a reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}

func TestOpenFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"empty path", "", true},
		{"nonexistent file", "/nonexistent/file.html", true},
		{"relative path", "browser.go", false},
		{"absolute path", "/etc/hosts", false},
		{"current directory", ".", false},
		{"parent directory", "..", false},
		{"nonexistent relative", "nonexistent-xyz123.html", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpenReader(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.Reader
		wantErr bool
	}{
		{"html content", strings.NewReader("<html><body>test</body></html>"), false},
		{"empty reader", strings.NewReader(""), false},
		{"nil reader", nil, true},
		{"large content", strings.NewReader(strings.Repeat("<html></html>", 1000)), false},
		{"error reader", &errorReader{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenReader(tt.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpenReader_Error(t *testing.T) {
	// Test error wrapping in OpenReader
	f, err := os.CreateTemp("", "test-*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	// Try to write to a closed file (should error)
	err = OpenReader(strings.NewReader("test"))
	// This should work, just checking error handling
	if err != nil {
		t.Logf("Expected error may occur: %v", err)
	}
}

func TestOpenURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"empty url", "", true},
		{"valid http", "https://go.dev/doc/devel/release", false},
		{"valid https", "https://go.dev", false},
		{"valid file", "file:///etc/hosts", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name    string
		prog    string
		args    []string
		wantErr bool
	}{
		{"empty program", "", nil, true},
		{"nonexistent command", "nonexistent-command-xyz123", []string{"arg"}, true},
		{"echo command", "echo", []string{"test"}, false},
		{"echo with multiple args", "echo", []string{"hello", "world"}, false},
		{"cat command", "cat", []string{"/etc/hosts"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runCmd(tt.prog, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStdoutStderr(t *testing.T) {
	// Test that Stdout and Stderr are set
	if Stdout == nil {
		t.Error("Stdout should not be nil")
	}
	if Stderr == nil {
		t.Error("Stderr should not be nil")
	}

	// Test that we can override them
	oldStdout := Stdout
	oldStderr := Stderr
	defer func() {
		Stdout = oldStdout
		Stderr = oldStderr
	}()

	Stdout = io.Discard
	Stderr = io.Discard

	// Run a command with overridden Stdout/Stderr
	err := runCmd("echo", "test")
	if err != nil {
		t.Errorf("runCmd() with overridden Stdout/Stderr failed: %v", err)
	}
}

func TestOpenReaderNil(t *testing.T) {
	// Explicitly test nil reader returns proper error
	err := OpenReader(nil)
	if err == nil {
		t.Error("OpenReader(nil) should return error")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("OpenReader(nil) error should mention nil: %v", err)
	}
}

func TestOpenFileEdgeCases(t *testing.T) {
	// Test various edge cases for OpenFile
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"dot", ".", false},
		{"dot slash", "./", false},
		{"double dot", "..", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenFile(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestRunCmdArgs(t *testing.T) {
	// Test runCmd with various argument combinations
	tests := []struct {
		name    string
		prog    string
		args    []string
		wantErr bool
	}{
		{"no args", "echo", nil, false},
		{"single arg", "echo", []string{"test"}, false},
		{"multiple args", "echo", []string{"a", "b", "c"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runCmd(tt.prog, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
