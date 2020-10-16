package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// A Handler type to store the template data
type templateHandler struct {
	once     sync.Once
	filename string
	tmpl     *template.Template
}

// ServeHTTP method to statisfy Handler Interface
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Compile template only once
	t.once.Do(func() {
		t.tmpl = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	// Execute writes data object to the template and writes the output to the io.writer passed
	t.tmpl.Execute(w, r)
}

func main() {
	// new room Handler instance
	roomInstance := newRoom()
	// Route Handler - Internal Router
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", roomInstance)
	// Running the room in a seperate goroutine
	go roomInstance.run()
	// Start Server
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Println("Started the chat room on https://localhost:3000")
		log.Fatal(err)
	}
}
