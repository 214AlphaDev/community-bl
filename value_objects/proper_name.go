package value_objects

import "errors"

type ProperName struct {
	firstName string
	lastName  string
}

func (n ProperName) FirstName() string {
	return n.firstName
}

func (n ProperName) LastName() string {
	return n.lastName
}

func NewProperName(firstName string, lastName string) (ProperName, error) {

	if len(firstName) == 0 {
		return ProperName{}, errors.New("first name must be a non empty string")
	}

	if len(lastName) == 0 {
		return ProperName{}, errors.New("last name must be a non empty string")
	}

	return ProperName{
		firstName: firstName,
		lastName:  lastName,
	}, nil
}
