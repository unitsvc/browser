package browser

import "strings"

func ExampleOpenFile() {
	// Open a local HTML file in the browser
	_ = OpenFile("index.html")
	// Output:
}

func ExampleOpenReader() {
	// Open HTML content from a reader
	html := strings.NewReader("<html><body>Hello, World!</body></html>")
	_ = OpenReader(html)
	// Output:
}

func ExampleOpenURL() {
	// Open a URL in the browser
	_ = OpenURL("https://go.dev")
	// Output:
}
