package main

import _ "embed"
import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed template/index.html
var template string

type Options struct {
	Port           int
	RootDir        string
	FileServerMode bool
}

func parseOptions() Options {
	port := flag.Int("p", 8080, "Port")
	dir := flag.String("d", ".", "Root directory")
	fileServerMode := flag.Bool("f", false, "Enable file server mode")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Println("Error resolving directory:", err)
		os.Exit(1)
	}

	return Options{
		Port:           *port,
		RootDir:        absDir,
		FileServerMode: *fileServerMode,
	}
}

func htmlTemplate(body string, title string) string {
	return fmt.Sprintf(template, title, body)
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
		return
	}
	contentType := detectContentType(filepath.Base(path))
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Write(fileContent)
}

func dirView(root, path string) string {
	title := "Index of " + path
	html := "<h1>" + title + "</h1>"
	rpath := filepath.Join(root, path)
	dir, err := os.ReadDir(rpath)
	if err != nil {
		return htmlTemplate("<p>Can't read directory "+rpath+"</p>", "Error")
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

func serve(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		methodStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
		fmt.Println(methodStyle.Render("GET"), r.URL.Path)

		rpath := filepath.Join(opts.RootDir, r.URL.Path)
		stat, err := os.Stat(rpath)
		if err != nil {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
			fmt.Println(errorStyle.Render("ERROR:"), rpath, "not found")
			http.NotFound(w, r)
			return
		}

		if stat.IsDir() {
			if !strings.HasSuffix(r.URL.Path, "/") {
				http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
				return
			}

			if !opts.FileServerMode {
				indexFile := filepath.Join(rpath, "index.html")
				_, err = os.ReadFile(indexFile)
				if err == nil {
					serveFile(indexFile, w)
					return
				}
			}

			w.Write([]byte(dirView(opts.RootDir, r.URL.Path)))
		} else {
			serveFile(rpath, w)
		}
	}
}

func renderUrl(url string, port string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	portStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)

	result := style.Render(url)

	if port != "" {
		result += portStyle.Render(":" + port)
	}

	return result
}

func main() {
	opts := parseOptions()

	servingPath := opts.RootDir
	homeDir, err := os.UserHomeDir()

	if err == nil {
		servingPath = strings.Replace(servingPath, homeDir, "~", 1)
	}

	fmt.Println(fmt.Sprintf("Serving %s on", servingPath), renderUrl("http://localhost", strconv.Itoa(opts.Port)))

	http.HandleFunc("/", serve(opts))
	err = http.ListenAndServe(":"+strconv.Itoa(opts.Port), nil)
	if err != nil {
		panic(err)
	}
}
