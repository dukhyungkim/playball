package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Game struct {
	players    map[string]*Player
	joinTicker *time.Ticker
	joinChan   chan *Player
	startChan  chan bool
	finishChan chan bool
	isPlaying  bool
}

func NewGame() *Game {
	g := &Game{
		players:    make(map[string]*Player),
		joinChan:   make(chan *Player, 1),
		startChan:  make(chan bool, 1),
		finishChan: make(chan bool, 1),
		joinTicker: time.NewTicker(time.Second),
	}

	g.joinTicker.Stop()
	go g.joinManager()
	go g.startManager()
	go g.finishManager()
	return g
}

func (g *Game) JoinPlayer(req *JoinRequest) error {
	if g.isPlaying {
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
			log.Println("start to wait new player")
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
	for {
		<-g.startChan
		log.Println("start new game")
		g.isPlaying = true
		g.Start()
		log.Println("finish game")
		g.finishChan <- true
	}
}

func (g *Game) Start() {
	baseball := NewBaseBall(4)
	log.Println("random number generated", baseball.answer)
	counter := 10
	for range time.NewTicker(time.Second).C {
		log.Printf("wait for %d sec to request guessing\n", counter)
		counter--
		if counter == 0 {
			break
		}
	}

	countOfPlayers := len(g.players)
	players := make([]*Player, 0, countOfPlayers)
	for _, player := range g.players {
		players = append(players, player)
	}

	guessNumbers := make([]int, countOfPlayers)
	wg := sync.WaitGroup{}
	for i, player := range players {
		wg.Add(1)
		go func(idx int, p *Player) {
			defer wg.Done()
			guessed, err := g.requestGuessing(p)
			if err != nil {
				log.Println(err)
			}
			guessNumbers[idx] = guessed
		}(i, player)
	}
	wg.Wait()

	result, err := baseball.compareToAnswer(guessNumbers[0])
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("guessed: %d, result: %+v\n", guessNumbers[0], result)

	for _, player := range players {
		wg.Add(1)
		go func(p *Player) {
			defer wg.Done()
			log.Printf("notifing results to %s\n", p.Name)
			resultInfo := []*ResultInfo{
				{
					Name:   p.Name,
					Number: guessNumbers[0],
					Result: result,
				},
			}
			err = g.notifyResults(p, resultInfo)
			if err != nil {
				log.Println(err)
			}
		}(player)
	}
	wg.Wait()
}

func (g *Game) Guess() {

}

type GuessRequest struct {
	Length       int `json:"length"`
	RemainChance int `json:"remain_chance"`
}

type GuessResponse struct {
	Number int `json:"number"`
}

func (g *Game) requestGuessing(p *Player) (int, error) {
	guessReq := GuessRequest{
		Length:       4,
		RemainChance: p.RemainChance,
	}

	b, _ := json.Marshal(guessReq)

	resp, err := sendPost(p.Address.String()+"/guess", bytes.NewBuffer(b))
	if err != nil {
		return 0, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	var guessResp GuessResponse
	err = json.NewDecoder(resp.Body).Decode(&guessResp)
	if err != nil {
		return 0, err
	}
	log.Println(guessResp)
	return guessResp.Number, nil
}

type ResultInfo struct {
	Name   string  `json:"name"`
	Number int     `json:"number"`
	Result *Result `json:"result"`
}

func (g *Game) notifyResults(p *Player, results []*ResultInfo) error {
	b, _ := json.Marshal(results)
	log.Println(string(b))

	resp, err := sendPost(p.Address.String()+"/notify_results", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (g *Game) finishManager() {
	for {
		<-g.finishChan
		log.Println("handle finish")
		g.isPlaying = false
		g.players = make(map[string]*Player)
	}
}

func sendPost(url string, data io.Reader) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, data)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
