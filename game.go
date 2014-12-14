package main

import (
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

type Game struct {
	players []Player
}

func NewGame(players []Player) (g *Game) {
	g = new(Game)
	g.players = players
	return g
}

func (g *Game) AddPlayer(p Player) {
	if len(g.players) >= 4 {
		log.Println("[-] Games can only handle 4 players at a time!")
		return
	}

	g.players = append(g.players, p)
}

func (g *Game) Winner() Player {
	countAlivePlayers := 0
	var alivePlayer Player
	for _, player := range g.players {
		if player.Alive() {
			countAlivePlayers++
			alivePlayer = player
		}
	}

	if countAlivePlayers == 1 {
		return alivePlayer
	}

	return nil
}

func (g *Game) ResetPlayers() {
	for _, player := range g.players {
		player.Reset()
	}
}

func (g *Game) play() {

	for {
		winner := g.playRound()
		winner.IncScore()
		waitOnInput()
		g.ResetPlayers()
		drawBoard()
	}
}

func (g *Game) playRound() Player {

	step := -1

	// learn starting moves
	for i, player := range g.players {
		g.broadcastMove(i, player.PrevMove(0))
	}

	for g.Winner() == nil {
		step++
		time.Sleep(50 * time.Millisecond)

		for playerIndex, player := range g.players {

			// Get the player move. Incase it's a quit or something.
			move := player.NextMove()

			if player.Alive() {
				g.broadcastMove(playerIndex, move)
				drawPlayer(player)

				if !g.validMove(player) {
					player.SetAlive(false)
					break
				}
			}
		}
		drawStats(step, g.players)
		termbox.Sync()

	}

	return g.Winner()

	return nil
}

func (g *Game) broadcastMove(p int, move Move) {
	for _, player := range g.players {
		player.RecordMove(p, move)
	}
}

func (g *Game) validMove(p Player) bool {
	lastMove := p.PrevMove(0)

	if lastMove.x <= 0 || lastMove.y <= 0 || lastMove.x >= ARENA_WIDTH-1 || lastMove.y >= ARENA_HEIGHT-1 {
		return false
	}

	for _, other := range g.players {
		otherMoves := other.Moves()
		if other == p {
			otherMoves = otherMoves[0 : len(otherMoves)-1]
		}

		for _, otherMove := range otherMoves {
			if otherMove.x == lastMove.x && otherMove.y == lastMove.y {
				return false
			}
		}
	}
	return true
}
