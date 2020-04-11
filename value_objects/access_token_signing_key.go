package value_objects

import (
	"errors"
	"reflect"
)

type AccessTokenSigningKey struct {
	bytes []byte
}

func (k AccessTokenSigningKey) Bytes() []byte {
	return k.bytes
}

func NewAccessTokenSigningKey(key []byte) (AccessTokenSigningKey, error) {

	if len(key) != 1024 {
		return AccessTokenSigningKey{}, errors.New("invalid signing key")
	}

	if reflect.DeepEqual(key, make([]byte, 1024)) {
		return AccessTokenSigningKey{}, errors.New("invalid signing key - slice of 0 bytes")
	}

	return AccessTokenSigningKey{
		bytes: key,
	}, nil
}
