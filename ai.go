package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"os"

	"github.com/nsf/termbox-go"
)

var AI_DISPLAY_ON = true

type AIPlayer struct {
	arena   [][]bool
	prev    Move
	display Display
}

func NewAI() (AI AIPlayer) {
	AI.arena = buildArena()
	return AI
}

func (ai *AIPlayer) ReadState(conn net.Conn) (state State) {
	dec := json.NewDecoder(conn)
	err := dec.Decode(&state)
	if err != nil {
		ai.display.Debug("Error: Lost connection with server." + err.Error())
		ai.WaitForInput()
		os.Exit(1)
	}
	return state
}

func (ai *AIPlayer) SendMove(conn net.Conn, m Move) {
	enc := json.NewEncoder(conn)
	err := enc.Encode(m)
	if err != nil {
		ai.display.Debug("Error: Lost connection with server." + err.Error())
		ai.WaitForInput()
		os.Exit(1)
	}
}

func (ai *AIPlayer) WaitForInput() {
	termbox.PollEvent()
}

func RunAI() {
	ai := NewAI()

	if AI_DISPLAY_ON {
		ai.display = NewDisplay()
		ai.display.DrawBoard()
	}

	conn, err := net.Dial("tcp", HOST+PORT)
	if err != nil {
		ai.display.Debug("Error: Failed to connect to server:" + HOST + PORT + "\nm" + err.Error())
		ai.WaitForInput()
		os.Exit(1)
	}

	for {
		state := ai.ReadState(conn)

		if state.Step == 0 {
			if AI_DISPLAY_ON {
				ai.display.Reset()
			}
			ai.Reset()
		}

		ai.UpdateArena(state)

		if AI_DISPLAY_ON {
			ai.display.UpdateState(state)
			ai.display.Sync()
		}

		if state.IsGameOver() {
			// The board will get reset on `if state.Step == 0`
			continue
		}

		if state.Players[state.PlayerIndex].Alive {
			m := ai.NextMove(state)
			ai.SendMove(conn, m)
		}

	}
}

func printMap(arena [][]bool) {
	for y := 0; y < ARENA_HEIGHT; y++ {
		out := ""
		for x := 0; x < ARENA_WIDTH; x++ {
			if arena[x][y] {
				out += "1"
			} else {
				out += "0"
			}
		}
		log.Println(out)
	}
}

func (ai *AIPlayer) UpdateArena(state State) {
	for _, player := range state.Players {
		ai.arena[player.X][player.Y] = true
	}
}

func (ai *AIPlayer) Reset() {
	ai.arena = buildArena()
}

func (ai *AIPlayer) isGoodDirection(d int) bool {
	move := ai.prev
	guessMove := UpdateMove(d, move)

	if guessMove.X <= 0 || guessMove.Y <= 0 || guessMove.X >= ARENA_WIDTH || guessMove.Y >= ARENA_HEIGHT {
		return false
	}

	return !ai.arena[guessMove.X][guessMove.Y]
}

func (ai *AIPlayer) NextMove(state State) Move {
	ai.prev = state.Players[state.PlayerIndex].Move

	nextMove := ai.prev
	// Continue on path
	if ai.isGoodDirection(nextMove.D) {
		nextMove = UpdateMove(nextMove.D, nextMove)
		ai.prev = nextMove
		return nextMove
	}

	possibleDirections := []int{UP, DOWN, LEFT, RIGHT}

	for len(possibleDirections) != 0 {
		i := rand.Intn(len(possibleDirections))
		direction := possibleDirections[i]
		if ai.isGoodDirection(direction) {
			nextMove = UpdateMove(direction, nextMove)
			ai.prev = nextMove
			return nextMove
		}
		possibleDirections = append(possibleDirections[0:i], possibleDirections[i+1:]...)
	}

	// there is no good move. :(
	nextMove = UpdateMove(nextMove.D, nextMove)
	ai.prev = nextMove
	return nextMove
}

func buildArena() [][]bool {
	arena := make([][]bool, ARENA_WIDTH)

	for i := 0; i < ARENA_WIDTH; i++ {
		col := make([]bool, ARENA_HEIGHT)
		col[0] = true
		col[ARENA_HEIGHT-1] = true
		arena[i] = col
	}

	for i := 0; i < ARENA_HEIGHT; i++ {
		arena[0][i] = true
		arena[ARENA_WIDTH-1][i] = true
	}
	return arena
}
