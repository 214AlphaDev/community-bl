package community_bl

import (
	"github.com/satori/go.uuid"
	vo "github.com/214alphadev/community-bl/value_objects"
	"time"
)

type ApplicationEntity struct {
	ID       ApplicationID
	MemberID MemberIdentifier
	ApplicationText string
	State ApplicationState
	RejectionReason string
	CreatedAt       time.Time
	RejectedAt      *time.Time
	ApprovedAt      *time.Time
	RejectedBy      *MemberIdentifier
	ApprovedBy      *MemberIdentifier
}

type MemberEntity struct {
	ID                    MemberIdentifier
	CreatedAt             time.Time
	VerifiedEmailAddress  bool
	Username              vo.Username
	EmailAddress          vo.EmailAddress
	Metadata              MetadataEntity
	MemberAccessPublicKey *vo.MemberAccessPublicKey
	AccessTokenID         *uuid.UUID
	Admin                 bool
	Verified bool
}

type MetadataEntity struct {
	ProperName   vo.ProperName
	ProfileImage *vo.Base64String
}

type MemberAccessTokenEntity struct {
	ExpiresAt         int64
	ID                uuid.UUID
	IssuedAt          int64
	Subject           MemberIdentifier
	signedAccessToken string
}

func (e MemberAccessTokenEntity) SignedAccessToken() string {
	return e.signedAccessToken
}

type ConfirmationCode struct {
	ID               uuid.UUID
	MemberIdentifier MemberIdentifier
	EmailAddress     vo.EmailAddress
	ConfirmationCode vo.ConfirmationCode
	IssuedAt         int64
	Used             bool
}

func (cc *ConfirmationCode) Expired() bool {
	return cc.IssuedAt+30*60 <= time.Now().Unix()
}
