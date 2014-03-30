package models

import (
	"github.com/albrow/zoom"
)

type SMSSubscriber struct {
	Phone string `zoom:"index"`
	zoom.DefaultData
}

func FindSMSSubscriberByNumber(num string) (*SMSSubscriber, error) {
	subs := make([]*SMSSubscriber, 0)
	q := zoom.NewQuery("SMSSubscriber")
	q = q.Filter("Phone =", num)
	if err := q.Scan(&subs); err != nil {
		return nil, err
	} else {
		if len(subs) == 0 {
			return nil, nil
		} else {
			return subs[0], nil
		}
	}
}
