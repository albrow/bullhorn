package controllers

import (
	"bullhorn/app/models"
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/albrow/zoom"
	"github.com/revel/revel"
	"regexp"
	"strings"
)

type Users struct {
	*revel.Controller
}

func (c Users) Index() revel.Result {
	users := make([]*models.User, 0)
	if err := zoom.NewQuery("User").Scan(&users); err != nil {
		c.RenderError(err)
	}
	return c.Render(users)
}

func (c Users) New() revel.Result {
	return c.Render()
}

func (c Users) Create(user *models.User) revel.Result {
	// strip extraneous characters from phone number
	user.Phone = formatPhoneNumber(user.Phone)

	validateUser(c, user, true)
	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Users.New)
	}

	crypted, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return c.RenderError(err)
	} else {
		user.EncryptedPass = string(crypted)
	}

	if err := zoom.Save(user); err != nil {
		c.RenderError(err)
	}
	url := "/users/" + user.Id
	return c.Redirect(url)
}

func (c Users) Edit(id string) revel.Result {
	user := new(models.User)
	if err := zoom.ScanById(id, user); err != nil {
		c.RenderError(err)
	}

	return c.Render(user)
}

func (c Users) Update(user *models.User) revel.Result {
	// strip extraneous characters from phone number
	user.Phone = formatPhoneNumber(user.Phone)

	validateUser(c, user, false)

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		url := "/users/" + user.Id + "/edit"
		return c.Redirect(url)
	}
	oldUser := new(models.User)
	if err := zoom.ScanById(user.Id, oldUser); err != nil {
		c.RenderError(err)
	}
	oldUser.Name = user.Name
	oldUser.Email = user.Email
	oldUser.Phone = user.Phone
	if err := zoom.Save(oldUser); err != nil {
		c.RenderError(err)
	}
	return c.Redirect(Users.Index)
}

func (c Users) Delete(id string) revel.Result {
	if err := zoom.DeleteById("User", id); err != nil {
		c.RenderError(err)
	}

	return c.Redirect(Users.Index)
}

func (c Users) SignIn() revel.Result {
	return c.RenderTemplate("Users/sign_in.html")
}

func (c Users) Authenticate(email, password string) revel.Result {
	c.Validation.Required(email)
	c.Validation.Required(password)
	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Users.SignIn)
	}

	user, err := models.FindUserByEmail(email)
	if err != nil {
		c.RenderError(err)
	} else if user == nil || bcrypt.CompareHashAndPassword([]byte(user.EncryptedPass), []byte(password)) != nil {
		c.Flash.Error("Invalid email/password combination")
		return c.Redirect(Users.SignIn)
	}

	c.Flash.Success("Signed In Successfully. Welcome, " + user.Name + "!")
	return c.Redirect("/")
}

func formatPhoneNumber(num string) string {
	numeric := "0123456789"
	formatted := ""
	// strip everything that's not a number
	for _, c := range num {
		if strings.Index(numeric, string(c)) != -1 {
			formatted += string(c)
		}
	}
	if len(formatted) == 10 {
		// assume in US and add country code 1
		formatted = "+1" + formatted
	} else {
		// add the + for the country code
		formatted = "+" + formatted
	}
	return formatted
}

func validateUser(c Users, user *models.User, requirePass bool) {
	if requirePass {
		c.Validation.Required(user.Pass)
		c.Validation.MinSize(user.Pass, 6)
	}
	c.Validation.Required(user.Email)
	c.Validation.Match(user.Email, regexp.MustCompile(".+\\@.+\\..+")).Message("Must be a valid email address")
	c.Validation.Required(user.Name)
	c.Validation.Required(user.Phone)
	c.Validation.MinSize(user.Phone, 12).Message("Must be full phone number with area code and country code.")
	c.Validation.MaxSize(user.Phone, 12).Message("Must be full phone number with area code and country code.")
}
