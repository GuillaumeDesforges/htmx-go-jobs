package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const port int = 8000

type Job struct {
	Id      string
	Command string
	Status  string
}

type Account struct {
	Jobs []Job
}

var account = Account{
	Jobs: []Job{
		{
			Id:      "1",
			Command: "echo 'Hello world'",
			Status:  "success",
		},
	},
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	requestLogJsonBytes, marshalError := json.Marshal(map[string]string{
		"method": r.Method,
		"path":   r.URL.Path,
	})
	if marshalError != nil {
		log.Panicln(marshalError)
	}
	log.Println(string(requestLogJsonBytes))
	tmpl := template.Must(template.ParseGlob("templates/*"))
	tmpl.ExecuteTemplate(w, "index", account)
}

func HandleJobs(w http.ResponseWriter, r *http.Request) {
	job := Job{
		Id:      "id",
		Command: "command",
		Status:  "status",
	}
	tmpl := template.Must(template.ParseGlob("templates/*"))
	tmpl.ExecuteTemplate(w, "job-row", job)
}

func main() {
	// static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// API routes
	http.HandleFunc("/api/jobs", HandleJobs)
	// HTML routes
	http.HandleFunc("/", HandleIndex)

	// start server
	address := fmt.Sprintf(":%d", port)
	http.ListenAndServe(address, nil)
}
