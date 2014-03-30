package models

import (
	"github.com/albrow/zoom"
)

func Init() error {
	if err := zoom.Register(&Subscriber{}); err != nil {
		return err
	}
	if err := zoom.Register(&User{}); err != nil {
		return err
	}
	return nil
}
