package models

import (
	"github.com/albrow/zoom"
)

type VoiceSubscriber struct {
	Phone string `zoom:"index"`
	zoom.DefaultData
}

func FindVoiceSubscriberByNumber(num string) (*VoiceSubscriber, error) {
	subs := make([]*VoiceSubscriber, 0)
	q := zoom.NewQuery("VoiceSubscriber")
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
