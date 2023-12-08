package main

import (
	"fmt"
	"github.com/google/uuid"
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
	Jobs: []Job{},
}

var jobUpdateChannel = make(chan Job)

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
		jobId := uuid.Must(uuid.NewRandom())
		job := Job{
			Id:      jobId.String(),
			Command: newJobCommand,
			Status:  "received",
		}
		account.Jobs = append(account.Jobs, job)
		go runJob(job)
		// return row
		tmpl := template.Must(template.ParseGlob("templates/*"))
		tmpl.ExecuteTemplate(w, "job-row", job)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func updateJobs() {
	// process incoming job updates
	for job := range jobUpdateChannel {
		fmt.Printf("Received %s\n", job)
		// find job in account
		for idx, accountJob := range account.Jobs {
			if accountJob.Id == job.Id {
				// update the job in the account
				accountJob.Status = job.Status
				account.Jobs[idx] = accountJob
				break
			}
		}
	}
}

func main() {
	// Static files
	fsStatic := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))
	// Logs
	fsLogs := http.FileServer(http.Dir("./logs"))
	http.Handle("/logs/", http.StripPrefix("/logs/", fsLogs))
	// API routes
	http.HandleFunc("/api/jobs", HandleJobs)
	// HTML routes
	http.HandleFunc("/", HandleIndex)

	// Update jobs in the database
	go updateJobs()

	// Start server
	address := fmt.Sprintf(":%d", port)
	http.ListenAndServe(address, nil)
}
