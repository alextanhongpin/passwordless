package passwordless

import (
	"context"
	"errors"
)

var ErrCodeExpired = errors.New("code expired")

type User struct {
	ID string
}

type (
	repository interface {
		Create(ctx context.Context, userID string) (*Code, error)
		Delete(ctx context.Context, code string) error
		Find(ctx context.Context, code string) (*Code, error)
	}
)

type Passwordless struct {
	repo repository
}

func (p *Passwordless) Start(ctx context.Context, userID string) (*Code, error) {
	// Instead of passing in email, we can delegate the user finding to the
	// client side, and just focus on the creation of the passwordless
	// mechanism.
	// user, err := p.users.WithEmail(ctx, email)
	// if err != nil {
	//         return nil, err
	// }

	code, err := p.repo.Create(ctx, userID)
	if err != nil {
		return nil, err
	}

	return code, nil
}

func (p *Passwordless) Authorize(ctx context.Context, code string) (*User, error) {
	// We can query the database with the conditional for expiry, but then
	// it makes it harder to find the token that actually expired and
	// return the relevant error messages. Also, we need to still delete
	// the unused code if it already expired. We can run a cron job to
	// periodically (passively) delete it, but it is easier to do it
	// actively.
	c, err := p.repo.Find(ctx, code)
	if err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		if err := p.repo.Delete(ctx, code); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := p.repo.Delete(ctx, code); err != nil {
		return nil, err
	}

	return &User{ID: c.UserID}, nil
}

func New(repo repository) *Passwordless {
	return &Passwordless{repo}
}
