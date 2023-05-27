package sendemail

import (
	"fmt"
	"log"
	"os"

	"net/smtp"
)

// emails list might be big, and frequency of api calls also might be big, so there is better to pass it by pointer
func TryToSendEmailsBtcUahPrice(emails_set_to_send *map[string]struct{}, btc_uah_price float64) {

	auth := smtp.PlainAuth("", os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com")

	for email := range *emails_set_to_send {
		msg := []byte(fmt.Sprintf("Subject: Bitcoin price\r\n"+
			"From: O. Fedorov's Software Engineering School Assignment <o.fedorov.genesis.assignment@gmail.com>\r\n"+
			"To: %v\r\n"+
			"Hi! Current BTC-UAH price is %.2f!\r\n", email, btc_uah_price))

		err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("EMAIL_USERNAME"), []string{email}, msg)

		if err != nil {
			log.Printf("Unable to send email to: %v. Reason: %v", email, err.Error())
		}
	}
}
