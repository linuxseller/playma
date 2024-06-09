package main
import (
    "log"
    "net/http"
    "io"
    "html"
    "html/template"
    "os"
    "path/filepath"
)
func home(w http.ResponseWriter, r *http.Request) {
    data := HomeData{Songs:music_files}
    tmpl, err := template.ParseFiles("./templates/home.html")
    if err != nil {
        log.Fatal("Error parsing home template")
    }
    tmpl.Execute(w, data)
}
func song(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Failed to parse form data", http.StatusBadRequest)
            return
        }
        filename := r.Form.Get("song")
        file, err := os.Open(html.UnescapeString(filename))
        if err != nil {
            http.Error(w, "File not found", http.StatusNotFound)
            return
        }

        // Set the content type based on the file extension
        contentType := "application/octet-stream"
        if filepath.Ext(filename) == ".mp3" {
            contentType = "audio/mpeg"
        } else if filepath.Ext(filename) == ".wav" {
            contentType = "audio/x-wav"
        }

        // Set the content type header
        w.Header().Set("Content-Type", contentType)

        // Copy the file to the response writer
        _, err = io.Copy(w, file)
        if err != nil {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
