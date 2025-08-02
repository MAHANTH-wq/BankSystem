package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	//create token for a specific username and duration
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)
	//Verify Token for verification of the token and will return the payload if success
	VerifyToken(token string) (*Payload, error)
}
