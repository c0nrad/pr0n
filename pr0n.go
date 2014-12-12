package main

import (
	"bytes"
	"fmt"

	"github.com/nsf/termbox-go"
)

const (
	WINDOW_WIDTH  = 80
	WINDOW_HEIGHT = 25
)

type Direction int
type Component int

const (
	UP = iota
	DOWN
	LEFT
	RIGHT

	BALL
	SHAFT
	HEAD
)

type Move struct {
	x, y int
	d    Direction
}

type Player struct {
	moves []Move
}

func (m *Move) String() string {
	return fmt.Sprintf("(%d, %d)", m.x, m.y)
}

func penisCharacter(d Direction, c Component) rune {
	switch c {
	case BALL:
		switch d {
		case UP:
			return rune('m')
		case DOWN:
			return rune('m')
		case LEFT, RIGHT:
			return rune('8')
		}
	case SHAFT:
		switch d {
		case LEFT, RIGHT:
			return rune('=')
		case UP, DOWN:
			return rune('|')
		}
	case HEAD:
		switch d {
		case UP:
			return rune('^')
		case DOWN:
			return rune('v')
		case LEFT:
			return rune('<')
		case RIGHT:
			return rune('>')
		}
	}
	return rune('.')
}

// i=0 would be the most recent move
func (p *Player) getMove(i int) Move {
	if i+1 > len(p.moves) {
		return Move{-1, -1, DOWN}
	}
	return p.moves[len(p.moves)-1-i]
}

func (p *Player) addMove(m Move) {
	p.moves = append(p.moves, m)
}

func drawPlayer(p Player) {
	head := p.getMove(0)
	shaft1 := p.getMove(1)
	shaft2 := p.getMove(2)
	balls := p.getMove(3)
	trail := p.getMove(4)
	termbox.SetCell(head.x, head.y, penisCharacter(head.d, HEAD), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(shaft1.x, shaft1.y, penisCharacter(shaft1.d, SHAFT), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(shaft2.x, shaft2.y, penisCharacter(shaft2.d, SHAFT), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(balls.x, balls.y, penisCharacter(balls.d, BALL), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(trail.x, trail.y, rune('.'), termbox.ColorDefault, termbox.ColorDefault)

}

func drawLine(line []byte, row int) {
	for x, c := range bytes.Runes(line) {
		termbox.SetCell(x, row, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func debug(line []byte) {
	for x, c := range bytes.Runes(line) {
		termbox.SetCell(x, WINDOW_HEIGHT, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func drawBoard() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	line := bytes.Repeat([]byte{'-'}, WINDOW_WIDTH)
	drawLine(line, 0)

	body := bytes.Repeat([]byte{' '}, WINDOW_WIDTH-2)
	body = append([]byte{'|'}, body...)
	body = append(body, byte('|'))
	for i := 0; i < WINDOW_HEIGHT-2; i++ {
		drawLine(body, i+1)
	}

	drawLine(line, WINDOW_HEIGHT-1)
	termbox.Sync()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	drawBoard()

	var player Player
	player.moves = []Move{Move{1, 1, RIGHT}, Move{2, 1, RIGHT}, Move{3, 1, RIGHT}, Move{4, 1, RIGHT}}
	drawPlayer(player)
	termbox.Sync()

loop:
	for {
		nextMove := player.getMove(0)
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowUp:
				nextMove.y -= 1
				nextMove.d = UP
				player.addMove(nextMove)
			case termbox.KeyArrowRight:
				nextMove.x += 1
				nextMove.d = RIGHT
				player.addMove(nextMove)
			case termbox.KeyArrowDown:
				nextMove.y += 1
				nextMove.d = DOWN
				player.addMove(nextMove)
			case termbox.KeyArrowLeft:
				nextMove.x -= 1
				nextMove.d = LEFT
				player.addMove(nextMove)
			}
		case termbox.EventResize:
			drawBoard()
		}

		debug([]byte(nextMove.String()))
		drawPlayer(player)
		termbox.Sync()

	}

}
