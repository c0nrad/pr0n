package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	STEP_WAIT_MS = 100
	ROOM_SIZE    = 2
)

type Room struct {
	Connections []net.Conn
	Name        string
}

type Game struct {
	Players []NetworkPlayer
	Step    int
}

type NetworkPlayer struct {
	Connection  net.Conn
	Alive       bool
	Moves       []Move
	MoveChannel chan Move
	Name        string
}

func HandleRoom(room Room) {
	game := NewGame(room.Connections)
	for {
		game.Play()
		time.Sleep(time.Second * 1)
		game.Reset()
	}
}

func (g *Game) Reset() {
	for i := range g.Players {
		g.Players[i].Alive = true
		g.Players[i].Moves = []Move{GetStartingSpot(i)}
	}
}

func NewGame(connections []net.Conn) Game {
	g := Game{}
	g.Players = make([]NetworkPlayer, len(connections))

	for i := range g.Players {
		g.Players[i].Connection = connections[i]
		g.Players[i].Alive = true
		g.Players[i].MoveChannel = make(chan Move, 100)
		g.Players[i].Moves = []Move{GetStartingSpot(i)}
		g.Players[i].Name = "Player " + strconv.Itoa(i)

		go MoveListener(g.Players[i].Connection, g.Players[i].MoveChannel)
	}
	return g
}

func MoveListener(conn net.Conn, moveChannel chan Move) {
	fmt.Println("MoveListner")
	conn.SetDeadline(time.Time{})
	dec := json.NewDecoder(conn)
	for {
		var m Move
		err := dec.Decode(&m)
		if err != nil {
			panic(err)
		}
		fmt.Println("We recieved a move from the connection. Player", m)
		moveChannel <- m
	}
}

func (g *Game) GetState(step int) (s State) {
	s.Step = step
	players := make([]PlayerState, len(g.Players))
	for i := range g.Players {
		players[i].Alive = g.Players[i].Alive
		players[i].Name = g.Players[i].Name

		//return last move, for dead this may be previous
		if players[i].Alive {
			players[i].Move = g.Players[i].Moves[step]
		} else {
			players[i].Move = g.Players[i].Moves[len(g.Players[i].Moves)-1]
		}
	}

	s.Players = players
	return s
}

func (g *Game) Winner() int {
	alivePlayers := 0
	lastAlive := -2
	for i, player := range g.Players {
		if player.Alive {
			alivePlayers++
			lastAlive = i
		}
	}
	if alivePlayers >= 2 {
		return -1
	}
	return lastAlive
}

func (g *Game) Play() {
	g.Step = 1
	fmt.Println(g.Winner(), "WINNTER")
	for g.Winner() == -1 { // GameOver

		prevState := g.GetState(g.Step - 1)
		fmt.Println("Prev State", prevState)

		for i, player := range g.Players {
			prevState.PlayerIndex = i
			out := StateToJSON(prevState)
			player.Connection.Write(out)
		}

		// wait for new moves
		time.Sleep(STEP_WAIT_MS * time.Millisecond)

		// guess next move
		for i := range g.Players {
			if g.Players[i].Alive {
				g.Players[i].Moves = append(g.Players[i].Moves, UpdateMove(g.Players[i].Moves[g.Step-1].D, g.Players[i].Moves[g.Step-1]))
			} else {
				continue
			}

			exit := false
			for {
				select {
				case move := <-g.Players[i].MoveChannel:
					fmt.Println("WE RECIEVED A MOVE", move)
					// Only take direction, incase of fowl play, CACAWWWW
					move = UpdateMove(move.D, g.Players[i].Moves[g.Step-1])
					g.Players[i].Moves[g.Step] = move
				default:
					exit = true
				}
				if exit {
					break
				}
			}

		}

		g.ValidateLastStep()
		g.Step++
	}

	// Game over, send last state
	prevState := g.GetState(g.Step - 1)
	fmt.Println("Prev State", prevState)

	for i, player := range g.Players {
		prevState.PlayerIndex = i
		out := StateToJSON(prevState)
		player.Connection.Write(out)
	}
}

func (g *Game) ValidateLastStep() {
	for i, player := range g.Players {
		if !player.Alive {
			continue
		}

		lastMove := player.Moves[g.Step]
		if lastMove.X <= 0 || lastMove.Y <= 0 || lastMove.X >= ARENA_WIDTH-1 || lastMove.Y >= ARENA_HEIGHT-1 {
			g.Players[i].Alive = false
		}

		for otherPlayerIndex, otherPlayer := range g.Players {
			for otherMoveIndex, otherMove := range otherPlayer.Moves {
				if lastMove.X == otherMove.X && lastMove.Y == otherMove.Y && !(otherMoveIndex == g.Step && otherPlayerIndex == i) {
					g.Players[i].Alive = false
				}
			}
		}
	}
}
