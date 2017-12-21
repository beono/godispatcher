package main

import (
	"errors"
	"fmt"
	"github.com/beono/godispatcher"
	"github.com/beono/godispatcher/examples/manager"
)

func validateUser(event *godispatcher.Event) error {
	if user, ok := event.Data.(*manager.User); ok {
		if user.Email == "" {
			return errors.New("email can't be empty")
		}
	}
	return nil
}

func main() {

	observer := godispatcher.New()
	observer.On(manager.UserUpdateBefore, validateUser)

	UserManager := manager.UserManager{
		Observer: observer,
	}

	newUser := manager.User{
		ID:    0,
		Email: "",
	}

	fmt.Println(UserManager.Update(newUser))
}
