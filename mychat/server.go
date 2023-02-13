package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
)

func userJoin(uc *UserControl, mc *MessagesControl, w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("u")

	greeting := uc.addUser(user, mc.lastMessageIdx())
	m := fmt.Sprintf("#%s %s", greeting, user)
	mc.addMessage(m)
}

func userSaid(mc *MessagesControl, w http.ResponseWriter, r *http.Request) {
	user,_ := url.QueryUnescape(r.URL.Query().Get("u"))
	message,_ := url.QueryUnescape(r.URL.Query().Get("b"))

	mc.addMessage(fmt.Sprintf("-%s: %s", user, message))
}

func userLeft(room string, uc *UserControl, mc *MessagesControl, w http.ResponseWriter, r *http.Request) {
	user,_ := url.QueryUnescape(r.URL.Query().Get("u"))

	uc.removeUser(user)
	mc.addMessage(fmt.Sprintf("#left %s", user))
}

func userListen(uc *UserControl, mc *MessagesControl, w http.ResponseWriter, r *http.Request) {
	user,_ := url.QueryUnescape(r.URL.Query().Get("u")) 
	userIsActive := false

	for {
		behind := 0
		actual := 0

		mc.waiting.waitWhile( func () bool { 
			if(!uc.getUserData(user, func (ud int) {behind = ud} )) {
				return false
			}

			userIsActive = true
			actual = mc.lastMessageIdx()
			return behind == actual
		} )

		if(userIsActive) {
			uc.setUserData(user, actual)
			mc.copyMessages(behind, actual, func (x string) { io.WriteString(w, fmt.Sprint(x, "\n")) })
		}

		break
	}
}

func server(room string, port string) {
	messageCondLock := sync.Mutex{}

	uc := UserControl { users:make(map[string]int, 100) }
	
	mc := MessagesControl{
		messages:[1000]string{},
		messagesIdx:0,
		messagesLock:sync.Mutex{},
		waiting:MessageWaitControl{messageCond:*sync.NewCond(&messageCondLock), waitingCC:atomic.Int32{}}}
	
	http.HandleFunc(fmt.Sprint("/", room, "/join/"), func (w http.ResponseWriter, r *http.Request) { userJoin(&uc, &mc, w, r) })
	http.HandleFunc(fmt.Sprint("/", room, "/leave/"), func (w http.ResponseWriter, r *http.Request) { userLeft(room, &uc, &mc, w, r) })
	http.HandleFunc(fmt.Sprint("/", room, "/say/"), func (w http.ResponseWriter, r *http.Request) { userSaid(&mc, w, r) })
	http.HandleFunc(fmt.Sprint("/", room, "/listen/"), func (w http.ResponseWriter, r *http.Request) { userListen(&uc, &mc, w, r) })

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}