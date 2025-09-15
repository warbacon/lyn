package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/gabriel-vasile/mimetype"
)

func logRequest(req *http.Request) {
	greenBold := color.New(color.FgGreen, color.Bold).Sprintf
	fmt.Println(greenBold(req.Method), req.URL.Path)
}

func customHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	localPath := "." + req.URL.Path

	file, err := os.Stat(localPath)
	if err != nil {
		color.Red("File %s not found", localPath)
		http.NotFound(w, req)
		return
	}

	var data []byte
	var contentType string

	if file.IsDir() {
		// Try to serve index.html
		indexPath := filepath.Join(localPath, "index.html")
		_, err := os.Stat(indexPath)
		if err == nil {
			fmt.Println("Serving " + indexPath)
			data, err = os.ReadFile(indexPath)
			if err != nil {
				color.Red("File %s cannot be read", indexPath)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			contentType = "text/html; charset=utf-8"
		} else {
			// Directory listing
			dir, err := os.ReadDir(localPath)
			if err != nil {
				color.Red("Cannot read directory %s", localPath)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			dirHtml := "<html><head><title>Lyn</title></head><body><h1>Directory listing for " + req.URL.Path + "</h1><ul>\n"

			// Add parent directory link if not at root
			if req.URL.Path != "/" {
				parentPath := req.URL.Path
				if strings.HasSuffix(parentPath, "/") && len(parentPath) > 1 {
					parentPath = parentPath[:len(parentPath)-1]
				}
				parentPath = filepath.Dir(parentPath)
				if !strings.HasSuffix(parentPath, "/") {
					parentPath += "/"
				}
				dirHtml += "<li><a href=\"" + parentPath + "\">../</a></li>\n"
			}

			for _, file := range dir {
				fileName := file.Name()
				if file.IsDir() {
					fileName += "/"
				}
				dirHtml += "<li><a href=\"" + filepath.Join(req.URL.Path, fileName) + "\">" + fileName + "</a></li>\n"
			}
			dirHtml += "</ul></body></html>\n"
			data = []byte(dirHtml)
			contentType = "text/html; charset=utf-8"
		}
	} else {
		// Serve individual file
		data, err = os.ReadFile(localPath)
		if err != nil {
			color.Red("File %s cannot be read", localPath)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		mtype, err := mimetype.DetectFile(filepath.Base(localPath))
		if err == nil {
			contentType = mtype.String()
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

func main() {
	port := "8080"
	http.HandleFunc("/", customHandler)
	fmt.Printf("Server running on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
