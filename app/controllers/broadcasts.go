package controllers

import (
	"bullhorn/app/models"
	"github.com/albrow/zoom"
	"github.com/revel/revel"
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

	if sms {
		revel.INFO.Println("SENDING...")
		revel.INFO.Println("MESSAGE: ", body)
		subscribers := make([]*models.Subscriber, 0)
		if err := zoom.NewQuery("Subscriber").Include("Phone").Scan(&subscribers); err != nil {
			return c.RenderError(err)
		}
		revel.INFO.Println(len(subscribers), "SUBSCRIBERS")
		for _, s := range subscribers {
			revel.INFO.Println("TO: ", s.Phone)
			_, _, err := twilio.SendSMS("+19103780902", s.Phone, body, "", "")
			if err != nil {
				return c.RenderError(err)
			}
		}
	}
	c.Flash.Success("Your Message Was Sent")
	return c.Redirect(Broadcasts.New)
}
