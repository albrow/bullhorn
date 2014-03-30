package controllers

import (
	"bullhorn/app/models"
	"github.com/albrow/zoom"
	"github.com/revel/revel"
	"github.com/sfreiberg/gotwilio"
)

type Broadcasts struct {
	*revel.Controller
}

func (c Broadcasts) New() revel.Result {
	return c.Render()
}

func (c Broadcasts) Create(body string, voice bool, sms bool) revel.Result {
	c.Validation.Required(body)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Broadcasts.New)
	}

	revel.INFO.Println("SENDING...")
	revel.INFO.Println("MESSAGE: ", body)

	// SMS
	smsSubs := make([]*models.SMSSubscriber, 0)
	if err := zoom.NewQuery("SMSSubscriber").Include("Phone").Scan(&smsSubs); err != nil {
		return c.RenderError(err)
	}
	revel.INFO.Println(len(smsSubs), "SUBSCRIBERS")
	for _, s := range smsSubs {
		revel.INFO.Println("TO: ", s.Phone)
		_, _, err := twilio.SendSMS("+19103780902", s.Phone, body, "", "")
		if err != nil {
			revel.ERROR.Println(err)
			return c.RenderError(err)
		}
	}

	// VOICE
	voiceSubs := make([]*models.VoiceSubscriber, 0)
	if err := zoom.NewQuery("VoiceSubscriber").Include("Phone").Scan(&voiceSubs); err != nil {
		return c.RenderError(err)
	}
	revel.INFO.Println(len(voiceSubs), "SUBSCRIBERS")
	for _, s := range voiceSubs {
		revel.INFO.Println("TO: ", s.Phone)
		callbackParams := gotwilio.NewCallbackParameters("http://107.170.105.233/say/" + body)
		_, _, err := twilio.CallWithUrlCallbacks("+19103780902", s.Phone, callbackParams)
		if err != nil {
			revel.ERROR.Println(err)
			return c.RenderError(err)
		}
	}

	c.Flash.Success("Your Message Was Sent")
	return c.Redirect(Broadcasts.New)
}

func (c Broadcasts) Say(message string) revel.Result {
	return c.Render(message)
}
