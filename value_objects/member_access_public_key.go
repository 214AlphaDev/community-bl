package value_objects

import (
	"errors"
	"golang.org/x/crypto/ed25519"
	"reflect"
)

type MemberAccessPublicKey struct {
	bytes []byte
}

func (k MemberAccessPublicKey) Key() ed25519.PublicKey {
	return k.bytes
}

func NewMemberAccessPublicKey(bytes []byte) (MemberAccessPublicKey, error) {

	if len(bytes) != 32 {
		return MemberAccessPublicKey{}, errors.New("invalid member access public key")
	}

	if reflect.DeepEqual(bytes, make([]byte, 32)) {
		return MemberAccessPublicKey{}, errors.New("invalid member access public key - nil slice")
	}

	return MemberAccessPublicKey{
		bytes: bytes,
	}, nil

}
