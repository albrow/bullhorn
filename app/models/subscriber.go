package models

import (
	"github.com/albrow/zoom"
)

type Subscriber struct {
	Phone string `zoom:"index"`
	zoom.DefaultData
}

func FindSubscriberByNumber(num string) (*Subscriber, error) {
	subs := make([]*Subscriber, 0)
	q := zoom.NewQuery("Subscriber")
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
