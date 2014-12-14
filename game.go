package main

import (
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

type Game struct {
	players      []Player
	alivePlayers []Player
	//board   Board
}

func (g *Game) IsGameOver() bool {
	return false
}

func (g *Game) play() Player {
	g.alivePlayers = g.players

	step := -1
	for {
		step++
		time.Sleep(50 * time.Millisecond)

		if len(g.alivePlayers) == 1 {

			return g.alivePlayers[0]
		}

		for playerIndex := 0; playerIndex < len(g.alivePlayers); playerIndex++ {
			player := g.alivePlayers[playerIndex]

			move := player.NextMove()
			g.broadcastMove(playerIndex, move)
			drawPlayer(player)

			if !g.validMove(player) {
				g.alivePlayers = append(g.alivePlayers[0:playerIndex], g.alivePlayers[playerIndex+1:]...)
				break
			}
		}
		log.Println(len(g.alivePlayers), len(g.players))
		drawStats(step, g.players)
		termbox.Sync()

	}

	return nil
}

func (g *Game) broadcastMove(p int, move Move) {
	for _, player := range g.alivePlayers {
		player.RecordMove(p, move)
	}
}

func (g *Game) validMove(p Player) bool {
	lastMove := p.PrevMove(0)

	if lastMove.x <= 0 || lastMove.y <= 0 || lastMove.x >= ARENA_WIDTH-1 || lastMove.y >= ARENA_HEIGHT-1 {
		log.Println("Hit a wall")
		return false
	}

	for _, other := range g.players {
		otherMoves := other.Moves()
		if other == p {
			otherMoves = otherMoves[0 : len(otherMoves)-1]
		}

		for _, otherMove := range otherMoves {
			if otherMove.x == lastMove.x && otherMove.y == lastMove.y {
				log.Println("Hit a player")
				return false
			}
		}
	}
	return true
}
