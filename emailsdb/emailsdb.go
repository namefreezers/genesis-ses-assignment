package emailsdb

import (
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
)

var emails_set map[string]struct{}

var underlying_writer_to_close *os.File
var Buf_writer *bufio.Writer

func is_email_valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

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
		if is_email_valid(email) {
			res[email] = struct{}{}
		}
	}

	return res
}

func Init(emails_file_path string) {

	emails_set = read_emails_from_file(emails_file_path)

	underlying_writer_to_close, err := os.OpenFile(emails_file_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Can't create file database : %v", emails_file_path)
		Buf_writer = nil
	} else {
		Buf_writer = bufio.NewWriter(underlying_writer_to_close)
	}
}

func AddEmail(email string) error {
	if !is_email_valid(email) {
		return fmt.Errorf("invalid email: %v", email) // we can't subsribe to invalid email, so return error here also
	}

	// check if email exists
	_, exists := emails_set[email]
	if exists {
		return fmt.Errorf("email is already subsribed: %v", email)
	}

	// add to in-memory db. Single entry is quite small (~50bytes), so we can store ~20M emails per each GB of memory.
	emails_set[email] = struct{}{}

	// add to file database to persist state between launches
	fmt.Fprintf(Buf_writer, "%v\n", email) // ignore write error here, because we need to return success, if this emails was not subscribed already (persistance of state between lauches is not the main goal)
	Buf_writer.Flush()                     // the same error ignoring here

	return nil
}

func GetCurrentEmailsSet() *map[string]struct{} {
	return &emails_set
}

func CloseFileDb() {
	if underlying_writer_to_close != nil {
		err := underlying_writer_to_close.Close()
		if err != nil {
			log.Printf("Can't close file database : %v", underlying_writer_to_close.Name())
		}
	}
}
