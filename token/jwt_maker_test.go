package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mahanth/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, tokenPayload, err := jwtMaker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, tokenPayload)

	payload, err := jwtMaker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, tokenPayload, err := jwtMaker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, tokenPayload)

	payload, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidTokenAlgoNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	require.NoError(t, err)

	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	responsePayload, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrInvalidToken.Error())
	require.Nil(t, responsePayload)

}
