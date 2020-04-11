package community_bl

import (
	vo "github.com/214alphadev/community-bl/value_objects"
)

type MemberRepository interface {
	FetchByID(memberID MemberIdentifier) (*MemberEntity, error)
	Save(member MemberEntity) error
	IsUsernameTaken(username vo.Username) (bool, error)
	IsEmailAddressTaken(emailAddress vo.EmailAddress) (bool, error)
	FetchByEmailAddress(emailAddress vo.EmailAddress) (*MemberEntity, error)
}

type ApplicationRepository interface {
	FetchLast(member MemberIdentifier) (*ApplicationEntity, error)
	Save(application ApplicationEntity) error
	FetchByID(applicationID ApplicationID) (*ApplicationEntity, error)
	FetchByQuery(query ApplicationsQuery) ([]ApplicationEntity, error)
}

type ConfirmationCodeRepository interface {
	Fetch(emailAddress vo.EmailAddress, confirmationCode vo.ConfirmationCode) (*ConfirmationCode, error)
	Save(cc *ConfirmationCode) error
	Last(emailAddress vo.EmailAddress) (*ConfirmationCode, error)
}

type MemberAccessPublicKeyRepository interface {
	AlreadyUsed(memberAccessPublicKey vo.MemberAccessPublicKey) (bool, error)
	Save(memberAccessPublicKey vo.MemberAccessPublicKey) error
}

type AccessTokenRepository interface {
	Save(accessToken *MemberAccessTokenEntity) error
}
