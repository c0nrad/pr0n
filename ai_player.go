package main

import "log"

type AIPlayer struct {
	name  string
	moves []Move
	arena [][]bool
	index int
}

func NewAIPlayer(name string, id int) (p *AIPlayer) {
	p = new(AIPlayer)
	p.moves = []Move{StartingSpots[id]}
	p.moves[0].d = UP
	p.name = name
	p.index = id

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
	p.arena = arena

	//printMap(arena)

	return p
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

func (p *AIPlayer) NextMove() Move {
	nextMove := p.PrevMove(0)

	direction := nextMove.d
	if p.isGoodDirection(direction) {
		nextMove = updateMove(direction, nextMove)
	} else if p.isGoodDirection(UP) {
		nextMove = updateMove(UP, nextMove)
	} else if p.isGoodDirection(DOWN) {
		nextMove = updateMove(DOWN, nextMove)
	} else if p.isGoodDirection(LEFT) {
		nextMove = updateMove(LEFT, nextMove)
	} else if p.isGoodDirection(RIGHT) {
		nextMove = updateMove(RIGHT, nextMove)
	}

	p.addMove(nextMove)
	return nextMove
}

func updateMove(d Direction, m Move) Move {
	m.d = d
	switch d {
	case UP:
		m.y -= 1
	case DOWN:
		m.y += 1
	case LEFT:
		m.x -= 1
	case RIGHT:
		m.x += 1
	}
	return m
}

func (p *AIPlayer) isGoodDirection(d Direction) bool {
	move := p.PrevMove(0)
	guessMove := updateMove(d, move)
	return !p.arena[guessMove.x][guessMove.y]
}

func (p *AIPlayer) String() string {
	out := "Player: "
	for _, m := range p.moves {
		out += m.String()
	}
	return out
}

func (p *AIPlayer) Name() string {
	return p.name
}

// i=0 would be the most recent move
func (p *AIPlayer) PrevMove(i int) Move {
	if i+1 > len(p.moves) {
		return Move{-1, -1, NONE}
	}
	return p.moves[len(p.moves)-1-i]
}

func (p *AIPlayer) addMove(m Move) {
	p.moves = append(p.moves, m)
}

func (p *AIPlayer) Moves() []Move {
	return p.moves
}

func (p *AIPlayer) RecordMove(player int, m Move) {
	p.arena[m.x][m.y] = true
}

func (p *AIPlayer) Index() int {
	return p.index
}
