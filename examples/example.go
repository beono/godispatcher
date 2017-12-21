package main

import (
	"fmt"
	"math/rand"

	"github.com/beono/godispatcher"
	"github.com/beono/godispatcher/examples/manager"
)

func setUserID(event *godispatcher.Event) error {
	if User, ok := event.Data.(*manager.User); ok {
		User.ID = rand.Uint64()
		fmt.Printf("new User ID is '%v'\n", User.ID)
	}
	return nil
}

func cleanCache(event *godispatcher.Event) error {
	if User, ok := event.Data.(manager.User); ok {
		fmt.Printf("cleaned cache for User '%v'\n", User.ID)
	}
	return nil
}

func main() {

	dispatcher := godispatcher.New()
	dispatcher.On(manager.UserUpdateBefore, setUserID)
	dispatcher.On(manager.UserUpdateAfter, cleanCache)

	UserManager := manager.UserManager{
		Observer: dispatcher,
	}

	newUser := manager.User{
		ID:    0,
		Email: "example@example.com",
	}

	fmt.Println(UserManager.Update(newUser))
}
