package main

import (
	"encoding/json"
	"math/rand"
)

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
)

type Move struct {
	X, Y, D int
}

type State struct {
	PlayerIndex int
	Step        int
	Players     []PlayerState
}

type PlayerState struct {
	Name string
	Move
	Alive bool
}

func StateToJSON(s State) []byte {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return b
}

func JSONToState(in []byte) (s State) {
	err := json.Unmarshal(in, &s)
	if err != nil {
		panic(err)
	}

	return s
}

var StartingSpots = []Move{
	Move{ARENA_WIDTH / 4, ARENA_HEIGHT / 2, RIGHT},
	Move{3 * ARENA_WIDTH / 4, ARENA_HEIGHT / 2, LEFT},
	Move{ARENA_WIDTH / 2, ARENA_HEIGHT / 4, DOWN},
	Move{ARENA_WIDTH / 2, 3 * ARENA_HEIGHT / 4, UP},
}

func GetRandomDirection() int {
	directions := []int{UP, RIGHT, DOWN, LEFT}
	return directions[rand.Intn(4)]
}
func GetStartingSpot(index int) Move {
	move := StartingSpots[index]
	move.D = GetRandomDirection()
	return move
}

func UpdateMove(d int, m Move) Move {
	m.D = d
	switch d {
	case UP:
		m.Y -= 1
	case DOWN:
		m.Y += 1
	case LEFT:
		m.X -= 1
	case RIGHT:
		m.X += 1
	}
	return m
}

func (s *State) IsGameOver() bool {
	aliveCount := 0
	for _, player := range s.Players {
		if player.Alive {
			aliveCount++
		}
	}
	return aliveCount <= 1
}
