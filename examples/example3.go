package main

import (
	"fmt"
	"math/rand"

	"github.com/beono/godispatcher"
	"github.com/beono/godispatcher/examples/manager"
)

func setUserID2(event *godispatcher.Event) error {
	if User, ok := event.Data.(*manager.User); ok {
		User.ID = rand.Uint64()
		fmt.Printf("new User ID is '%v'\n", User.ID)
	}
	return nil
}

func cleanCache2(event *godispatcher.Event) error {
	if User, ok := event.Data.(manager.User); ok {
		fmt.Printf("cleaned cache for User '%v'\n", User.ID)
	}
	return nil
}

func main() {

	dispatcher := godispatcher.New()
	dispatcher.On(manager.UserUpdateBefore, godispatcher.Listener{Callback: setUserID2, Priority: 1})
	dispatcher.On(manager.UserUpdateAfter, godispatcher.Listener{Callback: cleanCache2, Priority: 1})

	UserManager := manager.UserManager{
		Emitter: dispatcher,
	}

	newUser := manager.User{
		ID:    0,
		Email: "example@example.com",
	}

	fmt.Println(UserManager.Update(newUser))
}
