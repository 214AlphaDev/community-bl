package community_bl

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	vo "github.com/214alphadev/community-bl/value_objects"
	"reflect"
	"time"
)

type accessTokenService struct {
	signingKey            vo.AccessTokenSigningKey
	accessTokenRepository AccessTokenRepository
}

func (s *accessTokenService) New(member MemberEntity) (string, error) {

	if reflect.DeepEqual(member.ID, uuid.UUID{}) {
		return "", errors.New("invalid member id")
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 20).Unix(),
		Id:        uuid.NewV4().String(),
		IssuedAt:  time.Now().Unix(),
		Subject:   member.ID.String(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedAccessToken, err := accessToken.SignedString(s.signingKey.Bytes())
	if err != nil {
		return "", err
	}

	parsedAccessToken, err := s.Parse(signedAccessToken)
	if err != nil {
		return "", err
	}

	if err := s.accessTokenRepository.Save(&parsedAccessToken); err != nil {
		return "", err
	}

	return signedAccessToken, nil

}

func (s accessTokenService) Parse(accessToken string) (MemberAccessTokenEntity, error) {

	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return s.signingKey.Bytes(), nil
	})
	if err != nil {
		return MemberAccessTokenEntity{}, err
	}
	if !parsedAccessToken.Valid {
		return MemberAccessTokenEntity{}, errors.New("invalid access token")
	}

	claims, k := parsedAccessToken.Claims.(*jwt.StandardClaims)
	if !k {
		return MemberAccessTokenEntity{}, errors.New("got wrong claims")
	}

	subject, err := uuid.FromString(claims.Subject)
	if err != nil {
		return MemberAccessTokenEntity{}, err
	}

	accessTokenID, err := uuid.FromString(claims.Id)
	if err != nil {
		return MemberAccessTokenEntity{}, err
	}

	return MemberAccessTokenEntity{
		ExpiresAt:         claims.ExpiresAt,
		ID:                accessTokenID,
		IssuedAt:          claims.IssuedAt,
		Subject:           subject,
		signedAccessToken: accessToken,
	}, nil

}
