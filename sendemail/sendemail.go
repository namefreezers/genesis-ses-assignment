package sendemail

import (
	"fmt"
	"log"
	"os"
	"sync"

	"net/smtp"
)

func tryToSendEmail(email_to string, btc_uah_price float64, auth smtp.Auth, wg *sync.WaitGroup) {
	if wg != nil {
		// we send emails asyncronously, so `WaitGroup` is needed to wait for all email-sending-goroutines
		defer wg.Done()
	}

	msg := []byte(fmt.Sprintf("Subject: Bitcoin price\r\n"+
		"From: O. Fedorov's Software Engineering School Assignment <o.fedorov.genesis.assignment@gmail.com>\r\n"+
		"To: %v\r\n"+
		"Hi! Current BTC-UAH price is %.2f!\r\n", email_to, btc_uah_price))

	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("EMAIL_USERNAME"), []string{email_to}, msg)

	// we don't need to response an error from api if we couldn't send mail, so we drop this error
	if err != nil {
		log.Printf("Unable to send email to: %v. Reason: %v", email_to, err.Error())
	}
}

// emails set might be big, and frequency of api calls also might be big,
// so there is better to pass it by pointer
func TryToSendEmailsBtcUahPrice(emails_set_to_send *map[string]struct{}, btc_uah_price float64) {

	auth := smtp.PlainAuth("", os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com")

	// We will send them asyncronously to speed-up all the mailing. So we will wait for each goroutine to finish
	var wg sync.WaitGroup
	wg.Add(len(*emails_set_to_send))

	// send emails separately to have possibility to set header `To:` for each recipient individually.
	for email_to_send := range *emails_set_to_send {
		go tryToSendEmail(email_to_send, btc_uah_price, auth, &wg)
	}

	wg.Wait()
}
