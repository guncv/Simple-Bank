package token

import (
	"time"

	"github.com/guncv/Simple-Bank/util"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, role util.Role, duration time.Duration) (string, *Payload, error)

	// Verify Token
	VerifyToken(token string) (*Payload, error)
}
