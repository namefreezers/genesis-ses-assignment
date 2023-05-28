package emailsdb

import (
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
)

// In-memory db. Single entry is quite small (~50bytes), so we can store ~20M emails per each GB of memory.
// There is only `map` data type in go, so this is a `set` implementation via `map`
var emails_set map[string]struct{}

var underlying_writer_to_close *os.File // File will be closed upon process death, but jist in case we leave pointer to close it.
var buf_writer *bufio.Writer

func is_email_valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Used to read emails from file-db upon startup.
// We keep emails in `set`-like data structure, because we need to search (O(1)) and add (O(1))
func read_emails_from_file(emails_file_path string) map[string]struct{} {
	var res map[string]struct{} = make(map[string]struct{})

	opened_file, err := os.Open(emails_file_path)
	if err != nil {
		log.Println(err.Error())
		return res
	}

	defer opened_file.Close()

	scanner := bufio.NewScanner(opened_file)
	for scanner.Scan() {
		email := scanner.Text()
		// Keep only valid emails
		if is_email_valid(email) {
			res[email] = struct{}{}
		}
	}

	return res
}

// Used to init DB upon startup.
// 1) Prepare "in-memory-db"
// 2) Open file and store it to append new subscribed emails to file
func Init(emails_file_path string) {

	emails_set = read_emails_from_file(emails_file_path)

	underlying_writer_to_close, err := os.OpenFile(emails_file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Can't create file database : %v", emails_file_path)
		buf_writer = nil
	} else {
		buf_writer = bufio.NewWriter(underlying_writer_to_close)
	}
}

// Is called from `/api/subscribe` api endpoint
func AddEmail(email string) error {
	if !is_email_valid(email) {
		return fmt.Errorf("invalid email: %v", email) // we can't subsribe to invalid email, so return an error here
	}

	// Check if email exists in DB. If already exists - return error
	_, exists := emails_set[email]
	if exists {
		return fmt.Errorf("email is already subsribed: %v", email)
	}

	// add to in-memory db. Single entry is quite small (~50bytes), so we can store ~20M emails per each GB of memory.
	emails_set[email] = struct{}{}

	// If file was opened without error.
	if buf_writer != nil {
		// add to file database to persist state between launches
		fmt.Fprintf(buf_writer, "%v\n", email) // ignore write error here, because we need to return success, if this emails was not subscribed already (persistance of state between lauches is not the main goal)
		buf_writer.Flush()                     // the same ignoring of error here
	}

	return nil
}

// Is called from `/api/sendEmails` api endpoint
// In-memory-db may be quite big, so we pass it by pointer
func GetCurrentEmailsSet() *map[string]struct{} {
	return &emails_set
}

// Close file just in case (but it won't be called, because http server is infinite loop, and file will be closed by OS when process is killed)
func CloseFileDb() {
	if underlying_writer_to_close != nil { // if it was not opened.
		err := underlying_writer_to_close.Close()
		if err != nil {
			log.Printf("Can't close file database : %v", underlying_writer_to_close.Name())
		}
	}
}
