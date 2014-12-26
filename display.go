package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

type Component int

const (
	BALL = iota
	SHAFT
	HEAD
	TRAIL
)

type Display struct {
	Moves              [][]Move
	Names              []string
	Scores             []int
	PlayerRainbowIndex map[int]int
}

var Rainbow = []termbox.Attribute{termbox.ColorRed, termbox.ColorGreen, termbox.ColorYellow, termbox.ColorBlue, termbox.ColorMagenta, termbox.ColorCyan}

func NewDisplay() (d Display) {
	InitTermbox()
	SetupLogging()

	d.PlayerRainbowIndex = make(map[int]int)
	d.Moves = make([][]Move, ROOM_SIZE)
	d.Names = make([]string, ROOM_SIZE)
	d.Scores = make([]int, ROOM_SIZE)
	return d
}

func (d *Display) Reset() {
	d.PlayerRainbowIndex = make(map[int]int)
	d.Moves = make([][]Move, ROOM_SIZE)
	d.Names = make([]string, ROOM_SIZE)
	d.DrawBoard()
}

func InitTermbox() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.Sync()
}

func SetupLogging() {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)
}

func CloseTermbox() {
	termbox.Close()
}

func (d *Display) Sync() {
	termbox.Sync()
}

func PenisCharacter(d int, c Component) rune {
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

func (d *Display) DrawLine(line []byte, row int) {
	for x, c := range bytes.Runes(line) {
		termbox.SetCell(x, row, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (d *Display) DrawLineAt(line []byte, row, col int) {
	for x, c := range bytes.Runes(line) {
		termbox.SetCell(x+col, row, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (d *Display) Debug(line string) {
	for x, c := range bytes.Runes([]byte(line)) {
		termbox.SetCell(x, ARENA_HEIGHT, c, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

func (d Display) DrawBoard() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	line := bytes.Repeat([]byte{'-'}, ARENA_WIDTH)
	d.DrawLine(line, 0)

	body := bytes.Repeat([]byte{' '}, ARENA_WIDTH-2)
	body = append([]byte{'|'}, body...)
	body = append(body, byte('|'))
	for i := 0; i < ARENA_HEIGHT-2; i++ {
		d.DrawLine(body, i+1)
	}

	d.DrawLine(line, ARENA_HEIGHT-1)
}

func PrevMove(offset int, moves []Move) Move {
	if offset >= len(moves) {
		out := moves[0]
		out.X = -1
		out.Y = -1
		return out
	}

	return moves[len(moves)-offset-1]
}

func (d *Display) UpdateState(state State) {
	for i, player := range state.Players {
		d.Names[i] = player.Name
		d.Scores[i] = player.Score
		if player.Alive {
			d.Moves[i] = append(d.Moves[i], player.Move)
		} else {
			// if dead, at least make sure dead position is displayed
			if d.Moves[i][len(d.Moves[i])-1] != player.Move {
				d.Moves[i] = append(d.Moves[i], player.Move)
			}
		}
		d.DrawPlayer(i, d.Moves[i], state.Step)
	}

	d.DrawLineAt([]byte("Step: "+strconv.Itoa(state.Step)), 1, ARENA_WIDTH+1)

}

func (d *Display) DrawPlayer(playerIndex int, moves []Move, step int) {
	rIndex, ok := d.PlayerRainbowIndex[playerIndex]
	if !ok {
		rIndex = 0
	}

	head := PrevMove(0, moves)
	shaft1 := PrevMove(1, moves)
	shaft2 := PrevMove(2, moves)
	balls := PrevMove(3, moves)
	trail := PrevMove(4, moves)

	fmt.Println("head", head)
	termbox.SetCell(head.X, head.Y, PenisCharacter(head.D, HEAD), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(shaft1.X, shaft1.Y, PenisCharacter(shaft1.D, SHAFT), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(shaft2.X, shaft2.Y, PenisCharacter(shaft2.D, SHAFT), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(balls.X, balls.Y, PenisCharacter(balls.D, BALL), termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(trail.X, trail.Y, PenisCharacter(trail.D, TRAIL), Rainbow[rIndex], termbox.ColorDefault)

	rIndex++
	rIndex %= len(Rainbow)
	d.PlayerRainbowIndex[playerIndex] = rIndex

	d.DrawPlayerStats(playerIndex)
}

func lFill(in string, l int) string {
	return in + strings.Repeat(" ", l-len(in))
}

func (d *Display) DrawPlayerStats(playerIndex int) {
	moves := d.Moves[playerIndex]
	x := ARENA_WIDTH + 1
	y := 3 + playerIndex*ARENA_HEIGHT/4
	prevMove := PrevMove(0, moves)
	d.DrawLineAt([]byte(d.Names[playerIndex]+" ("+strconv.Itoa(d.Scores[playerIndex])+")"), y, x)
	d.DrawLineAt([]byte(fmt.Sprintf("(%d, %d)", prevMove.X, prevMove.Y)), y+1, x)
}
