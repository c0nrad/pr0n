package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nsf/termbox-go"
)

type Player struct {
	InputBuffer chan termbox.Event
	Prev        Move
}

const (
	PLAYER_DISPLAY_ON = true
)

func NewPlayer() (p Player) {
	termbox.Flush()
	p.InputBuffer = make(chan termbox.Event, 100)
	go KeyboardListener(p.InputBuffer)
	return
}

func KeyboardListener(input chan termbox.Event) {
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey && event.Key == termbox.KeyEsc {
			os.Exit(0)
		}
		input <- event
	}
	panic("unreachable")
}

func (p *Player) NextMove() Move {
	nextMove := Move{-1, -1, -1}

	for {
		select {
		case ev := <-p.InputBuffer:
			log.Println(ev)
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc, 'q', 'Q':
					os.Exit(0)
				case termbox.KeyArrowUp:
					if p.Prev.D == DOWN {
						continue
					}
					nextMove.D = UP
				case termbox.KeyArrowRight:
					if p.Prev.D == LEFT {
						continue
					}
					nextMove.D = RIGHT
				case termbox.KeyArrowDown:
					if p.Prev.D == UP {
						continue
					}
					nextMove.D = DOWN
				case termbox.KeyArrowLeft:
					if p.Prev.D == RIGHT {
						continue
					}
					nextMove.D = LEFT
				}
				if nextMove.D != p.Prev.D { // only eat one event?
					return nextMove
				}
			}
		default:
			return nextMove
		}
	}

	return nextMove
}

func ReadState(conn net.Conn) (state State) {
	dec := json.NewDecoder(conn)
	err := dec.Decode(&state)
	if err != nil {
		panic(err)
	}
	return state
}

func Play() {
	conn, err := net.Dial("tcp", HOST+PORT)
	if err != nil {
		panic(err)
	}

	var d Display
	if PLAYER_DISPLAY_ON {
		d = NewDisplay()
		d.DrawBoard()
	}

	p := NewPlayer()

	for {
		state := ReadState(conn)

		if state.Step == 0 {
			d.Reset()
		}

		if PLAYER_DISPLAY_ON {
			d.UpdateState(state)
			d.Sync()
		} else {
			fmt.Println("Recieved state", state)
		}

		if state.IsGameOver() {
			continue
		}

		p.Prev = state.Players[state.PlayerIndex].Move

		if state.Players[state.PlayerIndex].Alive {
			// If I am alive, send then ext move
			nextDirection := p.NextMove()
			if nextDirection.D != -1 { // we actually got a move...
				nextMove := UpdateMove(nextDirection.D, p.Prev)
				enc := json.NewEncoder(conn)
				err = enc.Encode(nextMove)
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
