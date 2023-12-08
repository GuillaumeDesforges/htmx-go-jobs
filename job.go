package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runJob(job Job) {
	job.Status = "starting"
	jobUpdateChannel <- job
	fmt.Printf("Start job id=%s\n", job.Id)

	job.Status = "running"
	jobUpdateChannel <- job
	jobOut, errCommand := exec.Command("bash", "-c", job.Command).Output()

	// update job status
	if errCommand != nil {
		job.Status = "failed"
		jobUpdateChannel <- job
	} else {
		job.Status = "done"
		jobUpdateChannel <- job
	}

	// write output to file
	logsFilePath := filepath.Join(".", "logs")
	errMkdir := os.MkdirAll(logsFilePath, os.ModePerm)
	if errMkdir != nil {
		log.Println(errMkdir)
	}
	jobOutFilePath := fmt.Sprintf("logs/%s.log", job.Id)
	errWriteLog := os.WriteFile(jobOutFilePath, jobOut, 0644)
	if errWriteLog != nil {
		log.Println(errWriteLog)
	}

	fmt.Printf("End job id=%s\n", job.Id)
}
