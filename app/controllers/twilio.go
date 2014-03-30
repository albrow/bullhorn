package controllers

import (
	"bullhorn/app/models"
	"github.com/albrow/zoom"
	"github.com/revel/revel"
)

type Twilio struct {
	*revel.Controller
}

func (c Twilio) SubscribeSMS() revel.Result {
	// bind variables to form body values
	var from, body string
	c.Params.Bind(&from, "From")
	c.Params.Bind(&body, "Body")
	revel.INFO.Println("RECEIVED TEXT")
	revel.INFO.Println("FROM: ", from)
	revel.INFO.Println("BODY: ", body)

	sub, err := models.FindSMSSubscriberByNumber(from)
	if err != nil {
		return c.RenderError(err)
	}
	if sub == nil {
		// number not yet registered
		switch body {
		case "start", "Start", "START":
			newSub := new(models.SMSSubscriber)
			newSub.Phone = from
			zoom.Save(newSub)
			return c.RenderText("You will now receive texts about events and services from UMD! Send STOP at any time to stop receiving texts.")
		case "stop", "Stop", "STOP":
			return c.RenderText("You will not receive texts about events and services from UMD. If you would like to, reply START.")
		default:
			return c.RenderText("Would you like to receive texts about events and services from UMD? Reply with START to confirm.")
		}
	} else {
		// number is already registered
		switch body {
		case "start", "Start", "START":
			return c.RenderText("You are already subscribed to receive texts about events and services from UMD!")
		case "stop", "Stop", "STOP":
			zoom.Delete(sub)
			return c.RenderText("You will no longer receive texts about events and services from UMD.")
		default:
			return c.RenderText("Would you like to stop receiving texts from us? Reply with STOP to stop.")
		}
	}
}

func (c Twilio) ReceiveVoice() revel.Result {
	var from string
	c.Params.Bind(&from, "From")
	revel.INFO.Println("RECEIVED CALL")
	revel.INFO.Println("FROM: ", from)
	sub, err := models.FindVoiceSubscriberByNumber(from)
	if err != nil {
		revel.ERROR.Println(err)
		return c.RenderError(err)
	}
	if sub == nil {
		return c.RenderTemplate("Twilio/receive-voice.xml")
	} else {
		return c.RenderTemplate("Twilio/already-subscribed-voice.xml")
	}
}

func (c Twilio) SubscribeVoice() revel.Result {
	var from, digits string
	c.Params.Bind(&from, "From")
	c.Params.Bind(&digits, "Digits")
	revel.INFO.Println("RECEIVED CALL")
	revel.INFO.Println("FROM: ", from)
	revel.INFO.Println("DIGITS: ", digits)

	sub, err := models.FindVoiceSubscriberByNumber(from)
	if err != nil {
		return c.RenderError(err)
	}
	if sub == nil {
		if digits == "1" {
			newSub := new(models.VoiceSubscriber)
			newSub.Phone = from
			if err := zoom.Save(newSub); err != nil {
				revel.ERROR.Println(err)
				return c.RenderError(err)
			}
			return c.RenderTemplate("Twilio/confirm-voice.xml")
		} else {
			return c.RenderTemplate("Twilio/confirm-mistake-voice.xml")
		}
	} else {
		if digits == "2" {
			zoom.Delete(sub)
			return c.RenderTemplate("Twilio/delete-voice.xml")
		} else {
			return c.RenderTemplate("Twilio/delete-mistake-voice.xml")
		}
	}
}
