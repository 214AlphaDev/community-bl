package community_bl

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	vo "github.com/214alphadev/community-bl/value_objects"
	"reflect"
	"time"
)

type memberService struct {
	onLogin []func(member MemberEntity)
	memberRepository                MemberRepository
	confirmationCodeRepository      ConfirmationCodeRepository
	transport                       Transport
	memberAccessPublicKeyRepository MemberAccessPublicKeyRepository
	accessTokenService              *accessTokenService
}

type RequestLoginCoolDownError struct {
	TryAgainAt int64
}

func (e RequestLoginCoolDownError) Error() string {
	return fmt.Sprintf("please retry to request the login at: %d", e.TryAgainAt)
}

func (s *memberService) SignUp(username vo.Username, emailAddress vo.EmailAddress, metadata MetadataEntity) (MemberEntity, error) {

	if reflect.DeepEqual(metadata, MetadataEntity{}) {
		return MemberEntity{}, errors.New("received empty metadata")
	}

	if reflect.DeepEqual(username, vo.Username{}) {
		return MemberEntity{}, errors.New("username value object was not correct initialized")
	}

	taken, err := s.memberRepository.IsUsernameTaken(username)
	if err != nil {
		return MemberEntity{}, err
	}
	if taken {
		return MemberEntity{}, errors.New("UsernameTaken")
	}

	if reflect.DeepEqual(emailAddress, vo.EmailAddress{}) {
		return MemberEntity{}, errors.New("email address value object was not correct initialized")
	}

	taken, err = s.memberRepository.IsEmailAddressTaken(emailAddress)
	if err != nil {
		return MemberEntity{}, err
	}
	if taken {
		return MemberEntity{}, errors.New("EmailAddressTaken")
	}

	member := MemberEntity{
		ID:           uuid.NewV4(),
		EmailAddress: emailAddress,
		Username:     username,
		Metadata:     metadata,
		CreatedAt:    time.Now(),
	}

	if err := s.memberRepository.Save(member); err != nil {
		return MemberEntity{}, err
	}

	return member, nil

}

func (s *memberService) RequestLogin(emailAddress vo.EmailAddress) error {

	member, err := s.memberRepository.FetchByEmailAddress(emailAddress)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("couldn't find member")
	}

	lastConfirmationCode, err := s.confirmationCodeRepository.Last(emailAddress)
	if err != nil {
		return err
	}

	if lastConfirmationCode != nil {
		if lastConfirmationCode.IssuedAt+120 >= time.Now().Unix() {
			return RequestLoginCoolDownError{
				TryAgainAt: time.Now().Unix() + ((lastConfirmationCode.IssuedAt + 120) - time.Now().Unix()),
			}
		}
	}

	code, err := vo.ConfirmationCodeFactory()
	if err != nil {
		return err
	}

	confirmationCode := &ConfirmationCode{
		ID:               uuid.NewV4(),
		EmailAddress:     emailAddress,
		ConfirmationCode: code,
		IssuedAt:         time.Now().Unix(),
		MemberIdentifier: member.ID,
	}
	if err := s.confirmationCodeRepository.Save(confirmationCode); err != nil {
		return err
	}

	return s.transport.SendConfirmationCode(*confirmationCode)

}

var LoginErrorConfirmationCodeNotFound = errors.New("confirmation code doesn't exist")
var LoginErrorConfirmationCodeExpired = errors.New("confirmation code expired")
var LoginErrorConfirmationCodeAlreadyUsed = errors.New("confirmation code already used")
var LoginErrorMemberAccessKeyHasAlreadyBeenUsed = errors.New("member access code has already been used")
var LoginErrorConfirmationCodeMemberMismatch = errors.New("member miss match - please try again")

func (s *memberService) Login(emailAddress vo.EmailAddress, memberAccessPublicKey vo.MemberAccessPublicKey, confirmationCode vo.ConfirmationCode) (MemberAccessTokenEntity, error) {

	if reflect.DeepEqual(emailAddress, vo.EmailAddress{}) {
		return MemberAccessTokenEntity{}, errors.New("email address is not a correctly initialized value object")
	}

	if reflect.DeepEqual(memberAccessPublicKey, vo.MemberAccessPublicKey{}) {
		return MemberAccessTokenEntity{}, errors.New("member access public key is not a correctly initialized value object")
	}

	if reflect.DeepEqual(confirmationCode, vo.ConfirmationCode{}) {
		return MemberAccessTokenEntity{}, errors.New("confirmation code is not a correctly initialized value object")
	}

	cc, err := s.confirmationCodeRepository.Fetch(emailAddress, confirmationCode)
	if err != nil {
		return MemberAccessTokenEntity{}, err
	}

	if cc == nil {
		return MemberAccessTokenEntity{}, LoginErrorConfirmationCodeNotFound
	}

	if cc.Expired() {
		return MemberAccessTokenEntity{}, LoginErrorConfirmationCodeExpired
	}

	if cc.Used {
		return MemberAccessTokenEntity{}, LoginErrorConfirmationCodeAlreadyUsed
	}

	used, err := s.memberAccessPublicKeyRepository.AlreadyUsed(memberAccessPublicKey)
	if err != nil {
		return MemberAccessTokenEntity{}, err
	}
	if used {
		return MemberAccessTokenEntity{}, LoginErrorMemberAccessKeyHasAlreadyBeenUsed
	}

	cc.Used = true
	if err := s.confirmationCodeRepository.Save(cc); err != nil {
		return MemberAccessTokenEntity{}, err
	}

	member, err := s.memberRepository.FetchByEmailAddress(emailAddress)
	if err != nil {
		return MemberAccessTokenEntity{}, nil
	}
	if member == nil {
		return MemberAccessTokenEntity{}, errors.New("couldn't find member")
	}

	if cc.MemberIdentifier != member.ID {
		return MemberAccessTokenEntity{}, LoginErrorConfirmationCodeMemberMismatch
	}

	signedAccessToken, err := s.accessTokenService.New(*member)
	if err != nil {
		return MemberAccessTokenEntity{}, err
	}

	accessToken, err := s.accessTokenService.Parse(signedAccessToken)
	if err != nil {
		return MemberAccessTokenEntity{}, err
	}

	member.MemberAccessPublicKey = &memberAccessPublicKey
	member.VerifiedEmailAddress = true
	member.AccessTokenID = &accessToken.ID
	if err := s.memberRepository.Save(*member); err != nil {
		return MemberAccessTokenEntity{}, err
	}
	if err := s.memberAccessPublicKeyRepository.Save(memberAccessPublicKey); err != nil {
		return MemberAccessTokenEntity{}, err
	}

	for _, onLogin := range s.onLogin {
		onLogin(*member)
	}

	return accessToken, nil

}

var GetMemberByAccessTokenErrorNoMember = errors.New("couldn't get member from access token")

func (s *memberService) GetByAccessToken(accessToken string) (MemberEntity, error) {

	fetchedAccessToken, err := s.accessTokenService.Parse(accessToken)
	if err != nil {
		return MemberEntity{}, err
	}

	member, err := s.memberRepository.FetchByID(fetchedAccessToken.Subject)
	if err != nil {
		return MemberEntity{}, err
	}

	if member == nil {
		return MemberEntity{}, GetMemberByAccessTokenErrorNoMember
	}

	return *member, nil

}

func (s memberService) GetMemberByID(id MemberIdentifier) (MemberEntity, error) {

	member, err := s.memberRepository.FetchByID(id)

	switch err {
	case nil:
		return *member, nil
	default:
		return MemberEntity{}, err
	}

}

func (s *memberService) OnLogin(cb func(member MemberEntity)) {
	s.onLogin = append(s.onLogin, cb)
}
