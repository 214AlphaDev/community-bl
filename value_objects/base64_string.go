package value_objects

import "encoding/base64"

type Base64String struct {
	data []byte
}

func (s Base64String) String() string {
	return base64.StdEncoding.EncodeToString(s.data)
}

func NewBase64String(str string) (Base64String, error) {

	data, err := base64.StdEncoding.DecodeString(str)
	switch err {
	case nil:
		return Base64String{
			data: data,
		}, nil
	default:
		return Base64String{}, err
	}

}
