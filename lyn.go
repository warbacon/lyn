package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func htmlTemplate(body string, title string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style> * { font-family: sans-serif; } html { color-scheme: light dark; } </style>
</head>
<body>
    %s
</body>
</html>
`, title, body)
}

func detectContentType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(filename, ".js"):
		return "text/javascript"
	case strings.HasSuffix(filename, ".css"):
		return "text/css"
	}

	return ""
}

func serveFile(path string, w http.ResponseWriter) {
	fileContent, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Can't read content of", path)
	}

	contentType := detectContentType(filepath.Base(path))
	w.Header().Add("Content-Type", contentType)
	w.Write(fileContent)
}

func dirView(root, path string) string {
	title := "Index of " + path
	html := "<h1>" + title + "</h1>"
	rpath := filepath.Join(root, path)
	dir, err := os.ReadDir(rpath)
	if err != nil {
		return htmlTemplate("<p>Can't read directory "+rpath+"</p>", "")
	}

	html += "\n<ul>"
	if path != "/" {
		parent := filepath.Dir(filepath.Clean(path))
		html += fmt.Sprintf("<li><a href=\"%s\">../</a></li>", parent)
	}
	for _, e := range dir {
		filename := e.Name()
		if e.IsDir() {
			filename += "/"
		}
		html += "\n<li><a href=\"" + path + filename + "\">" + filename + "</a></li>"
	}
	html += "\n</ul>"

	return htmlTemplate(html, title)
}

func serve(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET", r.URL.Path)
		rpath := filepath.Join(root, r.URL.Path)

		stat, err := os.Stat(rpath)
		if err != nil {
			fmt.Println("ERROR:", rpath, "not found")
			http.NotFound(w, r)
			return
		}

		if stat.IsDir() {
			if !strings.HasSuffix(r.URL.Path, "/") {
				http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
				return
			}

			w.Write([]byte(dirView(root, r.URL.Path)))
		} else {
			serveFile(rpath, w)
		}
	}
}

func main() {
	port := flag.Int("p", 8080, "Port")
	dir := flag.String("d", ".", "Root directory")
	flag.Parse()

	fmt.Printf("Serving HTTP on http://localhost:%d\n", *port)
	http.HandleFunc("/", serve(*dir))
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
