package main

import (
	"log"

	"github.com/namefreezers/genesis-ses-assignment/api"
	"github.com/namefreezers/genesis-ses-assignment/emailsdb"
)

const emails_file_path = "./emails_data/emails.txt"

func main() {
	log.SetPrefix("btc-course-service: ")

	emailsdb.Init(emails_file_path)
	defer emailsdb.CloseFileDb() // file will be released after death of the process, but just in case

	api.RunApi("0.0.0.0:5000")
}
