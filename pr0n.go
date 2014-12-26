package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

const (
	ARENA_WIDTH  = 80
	ARENA_HEIGHT = 25

	WINDOW_WIDTH  = 100
	WINDOW_HEIGHT = 30
)

var HOST = "localhost"
var PORT = ":1337"
var ROOM_SIZE = 2

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	flag.StringVar(&HOST, "host", HOST, "host to connect to")
	flag.StringVar(&PORT, "port", PORT, "port number to connect to")
	flag.IntVar(&ROOM_SIZE, "size", ROOM_SIZE, "number of players per room, max is 4")
	serverMode := flag.Bool("serve", false, "run in server move")
	aiMode := flag.Bool("ai", false, "run in server move")
	localMode := flag.Bool("local", false, "run everything locally")

	flag.Parse()

	if *serverMode {
		Server()
	} else if *aiMode {
		RunAI()
	} else if *localMode {
		AI_DISPLAY_ON = false
		SERVER_DEBUG = false
		go Server()
		fmt.Println("Starting server...")
		time.Sleep(50 * time.Millisecond)
		for i := 0; i < ROOM_SIZE-1; i++ {
			go RunAI()
		}
		Play()
	} else {
		Play()
	}
}
