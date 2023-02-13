package main

import (
	"fmt"
	"sync"
)

type UserControl struct {
	users map[string]int
	usersLock sync.Mutex
}

func withLocked[T any](l *sync.Mutex, op func () T) T {
	l.Lock()
	defer func() { l.Unlock() }()

	return op()
}

func (uc *UserControl) addUser(user string, lastMessage int) string {
	return withLocked( &uc.usersLock, func () string {
		_, exists := uc.users[user]
		if !exists {
			fmt.Println("addUser joined ", user)
			uc.users[user]=lastMessage // withLockedMessages(func () int { return messagesIdx} )
			return fmt.Sprintf("joined")
		}
		fmt.Println("addUser existed ", user)
		return fmt.Sprintf("returned")
	})
}

func (uc *UserControl) removeUser(user string) {
	withLocked( &uc.usersLock, func () int {
		fmt.Println("removeUser ", user)
		delete(uc.users, user)
		return 0
	})
}

func (uc *UserControl) getUserData(user string, op func (ud int)) bool {
	return withLocked( &uc.usersLock, func () bool {
		ud, exists := uc.users[user]
		if exists { 
			fmt.Println("getUserData ", user)
			op(ud) 
		}
		return exists
	})
}

func (uc *UserControl) setUserData(user string, ud int) bool {
	return withLocked( &uc.usersLock, func () bool {
		_, exists := uc.users[user]
		if exists {
			fmt.Println("setUserData ", user)
			uc.users[user] = ud
		}
		return exists
	})
}

