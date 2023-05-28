package main

import (
	"log"

	"github.com/namefreezers/genesis-ses-assignment/api"
	"github.com/namefreezers/genesis-ses-assignment/emailsdb"
)

const emails_file_path = "./emails_data/emails.txt"

func main() {
	log.SetPrefix("btc-course-service: ")

	// Read all emails from "file-db" to "in-memory-db" upon startup
	emailsdb.Init(emails_file_path)
	defer emailsdb.CloseFileDb() // file will be released after death of the process, but just in case

	// run gin http server
	api.RunApi("0.0.0.0:5000")
}
