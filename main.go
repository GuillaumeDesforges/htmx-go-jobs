package main

import (
	"fmt"
	"html/template"
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
	tmpl := template.Must(template.ParseGlob("templates/*"))
	tmpl.ExecuteTemplate(w, "index", account)
}

func HandleJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseGlob("templates/*"))
		for _, job := range account.Jobs {
			tmpl.ExecuteTemplate(w, "job-row", job)
		}
		return
	}
	if r.Method == http.MethodPost {
		// parse input
		r.ParseForm()
		newJobCommand := r.Form.Get("job-command")
		// create new job
		job := Job{
			Id:      "id",
			Command: newJobCommand,
			Status:  "queued",
		}
		account.Jobs = append(account.Jobs, job)
		// return row
		tmpl := template.Must(template.ParseGlob("templates/*"))
		tmpl.ExecuteTemplate(w, "job-row", job)
		return
	}
	w.WriteHeader(http.StatusNotFound)
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
