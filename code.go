package passwordless

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	CodeValidity = 10 * time.Minute
)

// Code represents a expirable, validatable code that is send to the end-user
// for verification.
type Code struct {
	UserID    string
	Code      string
	CreatedAt time.Time
	Validity  time.Duration
}

// Validate returns error if the code provided is expired or has been used.
func (c *Code) Validate() error {
	if time.Since(c.CreatedAt) > c.Validity {
		return ErrCodeExpired
	}
	return nil
}

type CodeOption func(c *Code)

func Validity(validity time.Duration) CodeOption {
	return func(c *Code) {
		c.Validity = validity
	}
}

func CustomCode(code string) CodeOption {
	return func(c *Code) {
		c.Code = code
	}
}

func NewCode(options ...CodeOption) *Code {
	c := &Code{
		Code:      uuid.Must(uuid.NewV4()).String(),
		CreatedAt: time.Now(),
		Validity:  CodeValidity,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}
