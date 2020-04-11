package value_objects

import (
	"fmt"
	"regexp"
)

var emailAddressRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type EmailAddress struct {
	emailAddress string
}

func (e EmailAddress) String() string {
	return e.emailAddress
}

func NewEmailAddress(emailAddress string) (EmailAddress, error) {

	if !emailAddressRegex.Match([]byte(emailAddress)) {
		return EmailAddress{}, fmt.Errorf("invalid email address: %s", emailAddress)
	}

	return EmailAddress{
		emailAddress: emailAddress,
	}, nil

}
