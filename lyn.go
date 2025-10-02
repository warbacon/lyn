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

func htmlTemplate(body string) string {
    return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lyn</title>
    <style> * { font-family: sans; } html { color-scheme: light dark; } </style>
</head>
<body>
    %s
</body>
</html>
`, body)
}

func detectMimeType(filename string) string {
    switch {
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

    mimetype := detectMimeType(filepath.Base(path))
    w.Header().Add("Content-Type", mimetype)
    w.Write(fileContent)
}

func dirView(root, path string) string {
    html := "<h1>Index of " + path + "</h1>"
    rpath := filepath.Join(root, path)
    dir, err := os.ReadDir(rpath)
    if err != nil {
        return htmlTemplate("<p>Can't read directory " + rpath + "</p>")
    }

    html += "\n<ul>"
    if path != "/" {
        html += "\n<a href=\"/\"><li>../</li></a>"
    }
    for _, e := range dir {
        filename := e.Name()
        if e.IsDir() {
            filename += "/"
        }
        html += "\n<a href=\"" + path + filename + "\"><li>" + filename + "</li></a>"
    }
    html += "\n</ul>"

    return htmlTemplate(html)
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

        if !stat.IsDir() {
            serveFile(rpath, w)
        } else {
            w.Write([]byte(dirView(root, r.URL.Path)))
        }
    }
}
func main() {
    port := flag.Int("p", 8080, "Port")
    dir := flag.String("d", ".", "Root directory")
    flag.Parse()
    fmt.Printf("Serving %s in http://localhost:%d\n", *dir, *port)
    http.HandleFunc("/", serve(*dir))
    http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
