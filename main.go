package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const port int = 8000

type Job struct {
	Id      string
	Command string
}

type Account struct {
	Jobs []Job
}

var account = Account{
	Jobs: []Job{
		{Id: "1", Command: "echo 'Hello world'"},
	},
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Print("Page index")
	tmpl := template.Must(template.ParseFiles("index.htmx"))
	tmpl.Execute(w, account)
}

func main() {
	// routes
	http.HandleFunc("/", Index)

	// start server
	address := fmt.Sprintf(":%d", port)
	http.ListenAndServe(address, nil)
}
