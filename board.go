package main

import (
	"bytes"
	"strings"

	"github.com/nsf/termbox-go"
)

type Board interface{}

type Direction int
type Component int

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
	case TRAIL:
		switch d {
		case UP, DOWN:
			return rune('|')
		case LEFT, RIGHT:
			return rune('-')
		}
	}
	return rune('.')
}

func drawLine(line []byte, row int) {
	for x, c := range bytes.Runes(line) {
		termbox.SetCell(x, row, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func drawLineAt(line []byte, row, col int) {
	for x, c := range bytes.Runes(line) {
		termbox.SetCell(x+col, row, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func debug(line string) {
	for x, c := range bytes.Runes([]byte(line)) {
		termbox.SetCell(x, ARENA_HEIGHT, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func drawBoard() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	line := bytes.Repeat([]byte{'-'}, ARENA_WIDTH)
	drawLine(line, 0)

	body := bytes.Repeat([]byte{' '}, ARENA_WIDTH-2)
	body = append([]byte{'|'}, body...)
	body = append(body, byte('|'))
	for i := 0; i < ARENA_HEIGHT-2; i++ {
		drawLine(body, i+1)
	}

	drawLine(line, ARENA_HEIGHT-1)
	termbox.Sync()
}

var playerRainbowIndex map[string]int
var Rainbow = []termbox.Attribute{termbox.ColorRed, termbox.ColorGreen, termbox.ColorYellow, termbox.ColorBlue, termbox.ColorMagenta, termbox.ColorCyan}

func init() {
	playerRainbowIndex = make(map[string]int)
}

func drawPlayer(p Player) {
	rIndex, ok := playerRainbowIndex[p.Name()]
	if !ok {
		rIndex = 0
	}

	head := p.PrevMove(0)
	shaft1 := p.PrevMove(1)
	shaft2 := p.PrevMove(2)
	balls := p.PrevMove(3)
	trail := p.PrevMove(4)
	termbox.SetCell(head.x, head.y, penisCharacter(head.d, HEAD), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(shaft1.x, shaft1.y, penisCharacter(shaft1.d, SHAFT), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(shaft2.x, shaft2.y, penisCharacter(shaft2.d, SHAFT), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(balls.x, balls.y, penisCharacter(balls.d, BALL), termbox.ColorDefault, termbox.ColorDefault)

	termbox.SetCell(trail.x, trail.y, penisCharacter(trail.d, TRAIL), Rainbow[rIndex], termbox.ColorDefault)

	rIndex++
	rIndex %= len(Rainbow)
	playerRainbowIndex[p.Name()] = rIndex
}

func lFill(in string, l int) string {
	return in + strings.Repeat(" ", l-len(in))
}

func drawStats(ts int, players []Player) {
	for _, player := range players {
		x := ARENA_WIDTH + 1
		y := 1 + player.Index()*WINDOW_HEIGHT/4
		prevMove := player.PrevMove(0)
		drawLineAt([]byte(lFill(player.Name(), WINDOW_WIDTH-ARENA_WIDTH)), y, x)
		drawLineAt([]byte(lFill(prevMove.String(), WINDOW_WIDTH-ARENA_WIDTH)), y+1, x)
	}

}
