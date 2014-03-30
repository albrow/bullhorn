package controllers

import (
	"github.com/revel/revel"
	"github.com/sfreiberg/gotwilio"
	"os"
)

var twilio *gotwilio.Twilio

func Init() {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	revel.INFO.Println(accountSid, authToken)
	twilio = gotwilio.NewTwilioClient(accountSid, authToken)
}
