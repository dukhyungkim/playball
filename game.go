package main

import (
	"log"
	"time"
)

type Game struct {
	players    map[string]*Player
	joinTicker *time.Ticker
	joinChan   chan *Player
	startChan  chan bool
	isStated   bool
}

func NewGame() *Game {
	g := &Game{
		players:    make(map[string]*Player),
		joinChan:   make(chan *Player, 1),
		startChan:  make(chan bool, 1),
		joinTicker: time.NewTicker(time.Second),
	}

	g.joinTicker.Stop()
	go g.joinManager()
	go g.startManager()
	return g
}

func (g *Game) JoinPlayer(req *JoinRequest) error {
	if g.isStated {
		return ErrGameStarted
	}

	_, has := g.players[req.Name]
	if has {
		return ErrDuplicateName
	}

	np, err := NewPlayer(req)
	if err != nil {
		return err
	}
	g.players[req.Name] = np
	g.joinChan <- np
	return nil
}

func (g *Game) joinManager() {
	const waitCounter = 10
	counter := waitCounter

	for {
		select {
		case <-g.joinChan:
			log.Println("reset join ticker")
			g.joinTicker.Reset(time.Second)
			counter = waitCounter
		case <-g.joinTicker.C:
			log.Printf("counter: %d\n", counter)
			counter--
			if counter == 0 {
				log.Println("fin waiting")
				g.joinTicker.Stop()
				g.startChan <- true
			}
		}
	}
}

func (g *Game) startManager() {
	<-g.startChan
	log.Println("start new game")
	g.isStated = true
}

func (g *Game) Start() {

}

func (g *Game) Guess() {

}
