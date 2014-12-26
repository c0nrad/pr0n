package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
)

const (
	AI_DISPLAY_ON = true
)

type AIPlayer struct {
	arena [][]bool
	prev  Move
}

func NewAI() (AI AIPlayer) {
	AI.arena = buildArena()
	return AI
}

func RunAI() {
	conn, err := net.Dial("tcp", HOST+PORT)
	if err != nil {
		panic(err)
	}

	AI := NewAI()

	var d Display
	if AI_DISPLAY_ON {
		d = NewDisplay()
		d.DrawBoard()
	}

	dec := json.NewDecoder(conn)

	for {
		var state State
		err := dec.Decode(&state)
		if err != nil {
			panic(err)
		}

		if state.Step == 0 {
			d.Reset()
			AI.Reset()
		}

		AI.UpdateArena(state)

		if AI_DISPLAY_ON {
			d.UpdateState(state)
			d.Sync()
		} else {
			fmt.Println("Recieved state", state)
		}

		if state.IsGameOver() {
			continue
		}

		if state.Players[state.PlayerIndex].Alive {
			// If I am alive, send then ext move
			m := AI.NextMove(state)
			enc := json.NewEncoder(conn)
			err = enc.Encode(m)
			if err != nil {
				panic(err)
			}
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
