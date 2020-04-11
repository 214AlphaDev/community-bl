package value_objects

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type ConfirmationCode struct {
	code string
}

func (c ConfirmationCode) String() string {
	return c.code
}

func NewConfirmationCode(code string) (ConfirmationCode, error) {

	if len(code) != 6 {
		return ConfirmationCode{}, fmt.Errorf("received invalid confirmation code with length: %d", len(code))
	}

	match, err := regexp.Match("^[0-9]*$", []byte(code))
	if err != nil {
		return ConfirmationCode{}, err
	}
	if !match {
		return ConfirmationCode{}, fmt.Errorf("confirmation code: %s is not a numeric string", code)
	}

	if code == "000000" {
		return ConfirmationCode{}, errors.New("invalid zero confirmation code")
	}

	return ConfirmationCode{
		code: code,
	}, nil

}

func ConfirmationCodeFactory() (ConfirmationCode, error) {

	const confirmationCodeL = 6

	var table = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	confCodeBytes := make([]byte, confirmationCodeL)
	readBytes, err := io.ReadAtLeast(rand.Reader, confCodeBytes, confirmationCodeL)
	if err != nil {
		return ConfirmationCode{}, err
	}
	if readBytes != confirmationCodeL {
		return ConfirmationCode{}, errors.New("failed to create confirmation code")
	}

	var confCode = ""
	for iteration := 0; iteration < len(confCodeBytes); iteration++ {
		confCode += strconv.Itoa(table[int(confCodeBytes[iteration])%len(table)])
	}

	return NewConfirmationCode(confCode)

}
