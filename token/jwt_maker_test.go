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

	role := util.DepositorRole
	token, tokenPayload, err := jwtMaker.CreateToken(username, role, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, tokenPayload)

	payload, err := jwtMaker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.Equal(t, payload.Role, role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	role := util.DepositorRole
	token, tokenPayload, err := jwtMaker.CreateToken(util.RandomOwner(), role, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, tokenPayload)

	payload, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidTokenAlgoNone(t *testing.T) {
	role := util.DepositorRole
	payload, err := NewPayload(util.RandomOwner(), role, time.Minute)
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
