package models

import (
	"fmt"
	"github.com/albrow/zoom"
)

type User struct {
	Email         string `zoom:"index"`
	Pass          string `redis:"-"`
	EncryptedPass string
	Name          string
	Phone         string `zoom:"index"`
	Admin         bool
	zoom.DefaultData
}

func (u *User) PrettyPhone() string {
	country_code := u.Phone[:2]
	area_code := u.Phone[2:5]
	three_digits := u.Phone[5:8]
	four_digits := u.Phone[8:]
	return fmt.Sprintf("%s (%s)-%s-%s", country_code, area_code, three_digits, four_digits)
}

func FindUserByEmail(email string) (*User, error) {
	users := make([]*User, 0)
	q := zoom.NewQuery("User")
	q = q.Filter("Email =", email)
	if err := q.Scan(&users); err != nil {
		return nil, err
	} else {
		if len(users) == 0 {
			return nil, nil
		} else {
			return users[0], nil
		}
	}
}
