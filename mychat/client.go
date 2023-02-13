package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func talkIsOver(roomTalkCh chan string, reason string) { 
	roomTalkCh <- reason
	close(roomTalkCh)
}

func sendUserTalksToRoom(roomUrl string, urlEscapedUser string, userCh chan string, roomTalkCh chan string) {

	defer talkIsOver(roomTalkCh, "#ended by sendUserTalksToRoom")

	for s := range userCh {
		switch resp, err := http.Get(fmt.Sprint(roomUrl, "say?u=", urlEscapedUser, "&b=", url.QueryEscape(s))); err {
		case nil: 
			body, err := io.ReadAll(resp.Body)	
			resp.Body.Close()
			if err != nil {
				roomTalkCh <- fmt.Sprint("!say response error: ", err)
				return
			}

			roomTalkCh <- string(body[:])
		default:
			roomTalkCh <- fmt.Sprint("!say error: ", err)
			return
		}
	}
}

func listenRoomTalks(roomUrl string, urlEscapedUser string, roomTalkCh chan string) {

	defer talkIsOver(roomTalkCh, "#ended by listenRoomTalks")

	for {
		switch resp, err := http.Get(fmt.Sprint(roomUrl, "listen?u=", urlEscapedUser)); err {
			case nil: 
				body, err := io.ReadAll(resp.Body)	
				resp.Body.Close()
				if err != nil {
					roomTalkCh <- fmt.Sprint("!listen response error: ", err)
					return
				}
				roomTalkCh <- string(body[:])
			default:
				roomTalkCh <- fmt.Sprint("!listen error: ", err)
				return
		}
	}
}

func listenUserTalks(serverCh chan string, roomTalkCh chan string) {
	defer talkIsOver(roomTalkCh, "#ended by listenUserTalks")

	stdInScanner := bufio.NewScanner(os.Stdin)
	for stdInScanner.Scan() {
		s := strings.ToValidUTF8(stdInScanner.Text(), "?")
		if(strings.Contains(s, "#leave")) {
			return
		}
		serverCh <- s
	}
}

func joinRoom(roomUrl string, urlEscapedUser string, roomOp func()) {
	resp, err := http.Get(fmt.Sprint(roomUrl, "join?u=", urlEscapedUser))
	if err != nil {
		fmt.Fprintln(os.Stderr, "!join error: ", err)
		return
	}
	defer func() { 
		cmd := fmt.Sprint(roomUrl, "leave?u=", urlEscapedUser)
		fmt.Fprintln(os.Stderr, cmd)
		http.Get(cmd) 
	}()

	body, err := io.ReadAll(resp.Body)	
	resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "!join read response error: ", err)
		return
	}

	fmt.Println(string(body[:]))

	roomOp()
}

func client(room string, addr string, user string) {

	userTalkCh := make(chan string)
	roomTalkCh := make(chan string)

	roomUrl := fmt.Sprint("http://", addr, "/", url.QueryEscape(room),"/")

	joinRoom(roomUrl, user, func() {
		fmt.Println("type #leave to leave the room", room)

		go sendUserTalksToRoom(roomUrl, url.QueryEscape(user), userTalkCh, roomTalkCh)
		go listenUserTalks(userTalkCh, roomTalkCh)
		go listenRoomTalks(roomUrl, url.QueryEscape(user), roomTalkCh)
		
		for ln := range roomTalkCh {
			fmt.Println(ln)
		}
	})
}