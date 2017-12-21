package manager

import (
	"github.com/beono/godispatcher"
	"github.com/pkg/errors"
)

// UserUpdateBefore event that happens before user updated
const UserUpdateBefore = "user.update.before"

// UserUpdateAfter event that happens after user updated
const UserUpdateAfter = "user.update.after"

// User is user entity
type User struct {
	ID       uint64
	Email    string
	IsActive bool
}

type emitter interface {
	On(event string, c godispatcher.Listener)
	Emit(event string, data interface{}) error
}

// UserManager talks to database
type UserManager struct {
	Emitter emitter
}

// Update updates a user in database
func (m UserManager) Update(u User) error {
	if err := m.Emitter.Emit(UserUpdateBefore, &u); err != nil {
		return errors.Wrap(err, "before update error")
	}

	// call database here

	if err := m.Emitter.Emit(UserUpdateAfter, u); err != nil {
		return errors.Wrap(err, "after update error")
	}
	return nil
}
