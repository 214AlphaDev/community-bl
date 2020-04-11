package community_bl

import (
	"errors"
	vo "github.com/214alphadev/community-bl/value_objects"
	"reflect"
)

type ApplicationsQuery struct {
	Position *ApplicationID
	Next     uint
	State    ApplicationState
}

type CommunityInterface interface {

	SignUp(username vo.Username, emailAddress vo.EmailAddress, metadata MetadataEntity) (MemberEntity, error)

	RequestLogin(emailAddress vo.EmailAddress) error

	Login(emailAddress vo.EmailAddress, memberAccessPublicKey vo.MemberAccessPublicKey, confirmationCode vo.ConfirmationCode) (MemberAccessTokenEntity, error)

	ApplyForVerification(applicationText string, member MemberIdentifier) (ApplicationEntity, error)

	ApproveApplication(applicationID ApplicationID, reviewer MemberIdentifier) error

	RejectApplication(applicationID ApplicationID, reason string, reviewer MemberIdentifier) error

	Applications(query ApplicationsQuery, requester MemberIdentifier) ([]ApplicationEntity, error)

	Application(application ApplicationID, requester MemberIdentifier) (ApplicationEntity, error)

	GetLastApplication(member MemberIdentifier, requester MemberIdentifier) (ApplicationEntity, error)

	GetMemberByAccessToken(accessToken string) (MemberEntity, error)

	GetMember(id MemberIdentifier) (MemberEntity, error)

	GetApplication(id ApplicationID) (ApplicationEntity, error)

	Promote(emailAddress vo.EmailAddress) error

	OnApplicationApproved(cb func(member MemberEntity))

	OnLogin(cb func(member MemberEntity))

}

type Community struct {
	communityService *communityService
	memberService    *memberService
}

func (c *Community) SignUp(username vo.Username, emailAddress vo.EmailAddress, metadata MetadataEntity) (MemberEntity, error) {
	return c.memberService.SignUp(username, emailAddress, metadata)
}

func (c *Community) RequestLogin(emailAddress vo.EmailAddress) error {
	return c.memberService.RequestLogin(emailAddress)
}

func (c *Community) Login(emailAddress vo.EmailAddress, memberAccessPublicKey vo.MemberAccessPublicKey, confirmationCode vo.ConfirmationCode) (MemberAccessTokenEntity, error) {
	return c.memberService.Login(emailAddress, memberAccessPublicKey, confirmationCode)
}

func (c *Community) ApplyForVerification(applicationText string, member MemberIdentifier) (ApplicationEntity, error) {
	return c.communityService.ApplyForVerification(member, applicationText)
}

func (c *Community) ApproveApplication(applicationID ApplicationID, reviewer MemberIdentifier) error {
	return c.communityService.ApproveApplication(applicationID, reviewer)
}

func (c *Community) RejectApplication(applicationID ApplicationID, reason string, reviewer MemberIdentifier) error {
	return c.communityService.RejectApplication(applicationID, reason, reviewer)
}

func (c *Community) Applications(query ApplicationsQuery, requester MemberIdentifier) ([]ApplicationEntity, error) {
	return c.communityService.Applications(query, requester)
}

func (c *Community) Application(application ApplicationID, requester MemberIdentifier) (ApplicationEntity, error) {
	return c.communityService.Application(application, requester)
}

func (c *Community) GetLastApplication(member MemberIdentifier, requester MemberIdentifier) (ApplicationEntity, error) {
	return c.communityService.GetLastApplication(member, requester)
}

func (c *Community) GetMemberByAccessToken(accessToken string) (MemberEntity, error) {
	return c.memberService.GetByAccessToken(accessToken)
}

func (c *Community) GetMember(id MemberIdentifier) (MemberEntity, error) {
	return c.memberService.GetMemberByID(id)
}

func (c *Community) GetApplication(id ApplicationID) (ApplicationEntity, error) {
	return c.communityService.GetApplicationByID(id)
}

func (c Community) Promote(emailAddress vo.EmailAddress) error {
	return c.communityService.Promote(emailAddress)
}

func (c *Community) OnApplicationApproved(cb func(member MemberEntity)) {
	c.communityService.OnApplicationApproved(cb)
}

func (c *Community) OnLogin(cb func(member MemberEntity)) {
	c.memberService.OnLogin(cb)
}

type Dependencies struct {
	MemberRepository                MemberRepository
	ApplicationRepository           ApplicationRepository
	ConfirmationCodeRepository      ConfirmationCodeRepository
	Transport                       Transport
	MemberAccessPublicKeyRepository MemberAccessPublicKeyRepository
	AccessTokenSigningKey           vo.AccessTokenSigningKey
	AccessTokenRepository           AccessTokenRepository
}

func NewCommunity(dependencies Dependencies) (*Community, error) {

	if reflect.DeepEqual(dependencies.AccessTokenSigningKey, vo.AccessTokenSigningKey{}) {
		return nil, errors.New("invalid access token signing key")
	}

	return &Community{
		communityService: &communityService{
			memberRepository:      dependencies.MemberRepository,
			applicationRepository: dependencies.ApplicationRepository,
		},
		memberService: &memberService{
			memberRepository:                dependencies.MemberRepository,
			confirmationCodeRepository:      dependencies.ConfirmationCodeRepository,
			transport:                       dependencies.Transport,
			memberAccessPublicKeyRepository: dependencies.MemberAccessPublicKeyRepository,
			accessTokenService: &accessTokenService{
				signingKey:            dependencies.AccessTokenSigningKey,
				accessTokenRepository: dependencies.AccessTokenRepository,
			},
		},
	}, nil

}
