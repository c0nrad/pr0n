package main

import (
	"fmt"
	"net"
)

func MatchMaker(newPlayers chan net.Conn) {
	currentRoom := Room{}

	for {
		select {
		case newPlayer := <-newPlayers:
			fmt.Println("Server: Match macker has recieved a new player!")
			currentRoom.Connections = append(currentRoom.Connections, newPlayer)

			if len(currentRoom.Connections) == ROOM_SIZE {
				fmt.Println("HERE?")
				go HandleRoom(currentRoom)
				currentRoom = Room{}
			}
		}
	}
}

func Server() {
	newPlayers := make(chan net.Conn, 100)
	go MatchMaker(newPlayers)

	ln, err := net.Listen("tcp", HOST+PORT)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		newPlayers <- conn
	}
}
