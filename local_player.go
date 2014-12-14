package main

type LocalPlayer struct {
	name  string
	moves []Move
	index int
	score int
	alive bool
}

func NewLocalPlayer(name string, id int) (p *LocalPlayer) {
	p = new(LocalPlayer)
	p.moves = []Move{StartingSpots[id]}
	p.name = name
	p.index = id
	p.alive = true
	return p
}

func (p *LocalPlayer) Reset() {
	p.alive = true
	p.moves = []Move{getStartingSpot(p.index)}
}

func (p *LocalPlayer) NextMove() Move {
	nextMove := p.PrevMove(0)
	direction := readInput().d
	if direction != NONE {
		nextMove.d = direction
	}
	switch nextMove.d {
	case UP:
		nextMove.y -= 1
	case DOWN:
		nextMove.y += 1
	case LEFT:
		nextMove.x -= 1
	case RIGHT:
		nextMove.x += 1
	}

	p.addMove(nextMove)
	return nextMove
}

func (p *LocalPlayer) String() string {
	out := "Player: "
	for _, m := range p.moves {
		out += m.String()
	}
	return out
}

func (p *LocalPlayer) Score() int {
	return p.score
}

func (p *LocalPlayer) IncScore() {
	p.score++
}

func (p *LocalPlayer) Alive() bool {
	return p.alive
}

func (p *LocalPlayer) SetAlive(alive bool) {
	p.alive = alive
}

func (p *LocalPlayer) Name() string {
	return p.name
}

// i=0 would be the most recent move
func (p *LocalPlayer) PrevMove(i int) Move {
	if i+1 > len(p.moves) {
		return Move{-1, -1, NONE}
	}
	return p.moves[len(p.moves)-1-i]
}

func (p *LocalPlayer) addMove(m Move) {
	p.moves = append(p.moves, m)
}

func (p *LocalPlayer) Moves() []Move {
	return p.moves
}

func (p *LocalPlayer) RecordMove(player int, m Move) {
	return
}

func (p *LocalPlayer) Index() int {
	return p.index
}
