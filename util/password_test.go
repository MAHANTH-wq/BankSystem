package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(6)
	hashedPassword1, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	// Check that the hashed password can be verified
	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	// Check that a wrong password fails verification
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.Error(t, err)
	require.Error(t, bcrypt.ErrMismatchedHashAndPassword, err)

	// Check that hashing the same password produces different hashes
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
