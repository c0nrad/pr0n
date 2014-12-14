package main

import (
	"log"
	"os"

	"github.com/nsf/termbox-go"
)

const (
	ARENA_WIDTH  = 80
	ARENA_HEIGHT = 25

	WINDOW_WIDTH  = 100
	WINDOW_HEIGHT = 30
)

var inputBuffer chan termbox.Event

func readInput() Move {
	nextMove := Move{0, 0, NONE}

	for {
		select {
		case ev := <-inputBuffer:

			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					os.Exit(0)
				case termbox.KeyArrowUp:
					nextMove.d = UP
				case termbox.KeyArrowRight:
					nextMove.d = RIGHT
				case termbox.KeyArrowDown:
					nextMove.d = DOWN
				case termbox.KeyArrowLeft:
					nextMove.d = LEFT
				}
			}
		default:
			return nextMove
		}
	}

	return nextMove
}

func waitOnInput() termbox.Event {
	done := false
	for {
		select {
		case <-inputBuffer:
			continue
		default:
			done = true
			break
		}

		if done {
			break
		}
	}

	select {
	case ev := <-inputBuffer:
		return ev
	}

}

func setupLogging() {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)
}

func main() {
	setupLogging()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	drawBoard()

	p1 := NewLocalPlayer("c0nrad", 0)
	p2 := NewAIPlayer("Sylvia", 1)
	p3 := NewAIPlayer("Thor", 2)
	p4 := NewAIPlayer("Bae", 3)

	termbox.Sync()

	inputBuffer = make(chan termbox.Event, 100)
	go func() {
		for {
			inputBuffer <- termbox.PollEvent()
		}
	}()

	pss := []Player{p1, p2, p3, p4}
	g := Game{pss, pss}

	waitOnInput()
	g.play()
	waitOnInput()

	debug(("WINNER"))
}
