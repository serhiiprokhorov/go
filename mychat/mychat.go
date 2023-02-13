package main

import (
	"fmt"
	"os"
)

func create_room_from_os_args(room_args []string) {
	if(len(room_args) < 4) {
		fmt.Println("room name is missing")
		return
	}
	server(room_args[2], room_args[4])
}

func join_room_from_os_args(room_args []string) {
	if(len(room_args) < 7) {
		fmt.Println("room name is missing")
		return
	}
	client(room_args[2], room_args[4], room_args[6])
}

func main() {
	// process args
	for i, s := range os.Args {
		if(s == "create") {
			create_room_from_os_args(os.Args[i:])
			return
		}
		if(s == "join") {
			join_room_from_os_args(os.Args[i:])
			return
		}
	}

	fmt.Println("mychat create room {room name} at {room port}")
	fmt.Println("mychat join room {room name} at {room address} as {user name} ")
	return
}