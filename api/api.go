package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namefreezers/genesis-ses-assignment/emailsdb"
	"github.com/namefreezers/genesis-ses-assignment/fetchbtcrate"
	"github.com/namefreezers/genesis-ses-assignment/sendemail"
)

func getRate(c *gin.Context) {
	btc_rate, err := fetchbtcrate.FetchBtcUahRateMain()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "all btc rate providers are unavailable"})
		return
	}

	c.IndentedJSON(http.StatusOK, btc_rate)
}

type subscribe_form_data struct {
	Email string `json:"email"`
}

func postSubscribe(c *gin.Context) {
	var form_data subscribe_form_data

	if err := c.BindJSON(&form_data); err != nil {
		// return 409 if wrong request
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "couldn't parse request args."})
		return
	}

	err := emailsdb.AddEmail(form_data.Email)
	if err != nil {
		// return 409 if 1) email is not valid or 2) email was already subscribed
		c.IndentedJSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func postSendEmails(c *gin.Context) {
	btc_uah_rate, err := fetchbtcrate.FetchBtcUahRateMain()
	if err != nil {
		log.Println(err.Error())
		c.IndentedJSON(http.StatusOK, gin.H{"message": "all btc rate providers are unavailable. Email weren't sent."})
		return
	}

	sendemail.TryToSendEmailsBtcUahPrice(emailsdb.GetCurrentEmailsSet(), btc_uah_rate)
	c.Status(http.StatusOK)
}

func RunApi(addr string) {

	router := gin.Default()

	router.GET("/api/rate", getRate)
	router.POST("/api/subscribe", postSubscribe)
	router.POST("/api/sendEmails", postSendEmails)

	router.Run(addr)
}
