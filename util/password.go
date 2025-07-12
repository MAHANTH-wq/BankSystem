package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// This function should implement the logic to hash the password.

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) error {
	// This function should implement the logic to check the password against the hashed password.

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
