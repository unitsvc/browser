package browser

import (
	"io"
	"strings"
	"testing"
)

func TestOpenFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"empty path", "", true},
		{"nonexistent file", "/nonexistent/file.html", true},
		{"relative path", "browser.go", false},
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

func TestOpenURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"empty url", "", true},
		{"valid http", "https://go.dev/doc/devel/release", false},
		{"valid https", "https://go.dev", false},
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
