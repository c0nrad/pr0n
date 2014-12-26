package main

import (
	"flag"
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	flag.StringVar(&HOST, "host", "c0nrad.io", "host to connect to")
	flag.StringVar(&PORT, "port", ":1337", "port number to connect to")
	serverMode := flag.Bool("serve", false, "run in server move")
	aiMode := flag.Bool("ai", false, "run in server move")

	flag.Parse()

	if *serverMode {
		Server()
	} else if *aiMode {
		RunAI()
	} else {
		PlayLocal()
	}
}
