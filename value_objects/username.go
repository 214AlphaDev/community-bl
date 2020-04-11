package value_objects

import (
	"fmt"
	"regexp"
)

type Username struct {
	username string
}

func (u Username) String() string {
	return u.username
}

func NewUsername(username string) (Username, error) {

	matched, err := regexp.Match("^_{0,1}([a-z0-9]+_{0,1})+$", []byte(username))
	if err != nil {
		return Username{}, err
	}

	if !matched {
		return Username{}, fmt.Errorf(`username: "%s" is invalid`, username)
	}

	return Username{
		username: username,
	}, nil
}
