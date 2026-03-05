// Package browser provides helpers to open files, readers, and urls in a browser window.
//
// The choice of which browser is started is entirely client dependent.
package browser

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Stdout is the io.Writer to which executed commands write standard output.
var Stdout io.Writer = os.Stdout

// Stderr is the io.Writer to which executed commands write standard error.
var Stderr io.Writer = os.Stderr

// OpenFile opens new browser window for the file path.
// The path is converted to an absolute path automatically.
func OpenFile(path string) error {
	if path == "" {
		return fmt.Errorf("browser: path cannot be empty")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("browser: failed to get absolute path: %w", err)
	}
	return OpenURL("file://" + absPath)
}

// OpenReader consumes the contents of r and presents the
// results in a new browser window.
func OpenReader(r io.Reader) error {
	if r == nil {
		return fmt.Errorf("browser: reader cannot be nil")
	}
	f, err := os.CreateTemp("", "browser.*.html")
	if err != nil {
		return fmt.Errorf("browser: could not create temporary file: %w", err)
	}
	defer func() {
		// Clean up temp file if we fail before OpenFile
		if err != nil {
			os.Remove(f.Name())
		}
		f.Close()
	}()

	if _, err = io.Copy(f, r); err != nil {
		return fmt.Errorf("browser: caching temporary file failed: %w", err)
	}
	return OpenFile(f.Name())
}

// OpenURL opens a new browser window pointing to url.
func OpenURL(url string) error {
	if url == "" {
		return fmt.Errorf("browser: URL cannot be empty")
	}
	return openBrowser(url)
}

func runCmd(prog string, args ...string) error {
	if prog == "" {
		return fmt.Errorf("browser: program cannot be empty")
	}
	cmd := exec.Command(prog, args...)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr
	return cmd.Run()
}
