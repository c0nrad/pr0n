package main

import "fmt"

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
	NONE

	BALL
	SHAFT
	HEAD
	TRAIL
)

type Move struct {
	x, y int
	d    Direction
}

func (m *Move) String() string {
	return fmt.Sprintf("(%d, %d, %d)", m.x, m.y, m.d)
}

type PlayerType int

const (
	LOCAL = iota
	AI
	NETWORK
)

var StartingSpots = []Move{
	Move{ARENA_WIDTH / 4, ARENA_HEIGHT / 2, RIGHT},
	Move{3 * ARENA_WIDTH / 4, ARENA_HEIGHT / 2, LEFT},
	Move{ARENA_WIDTH / 2, ARENA_HEIGHT / 4, DOWN},
	Move{ARENA_WIDTH / 2, 3 * ARENA_HEIGHT / 4, UP},
}

type Player interface {

	// Return the name of the player
	Name() string

	// Returns the players index
	Index() int

	// Return a list of all previous moves
	Moves() []Move

	// Get the next move from the player
	NextMove() Move

	// A convinience function, get the ith previous move
	PrevMove(i int) Move

	// Tell the player about a previous move
	RecordMove(player int, m Move)
}
