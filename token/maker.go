package token

import "time"

type Maker interface {
	// gen token for specific user
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// check token & get payload from it
	VerifyToken(token string) (*Payload, error)
}
