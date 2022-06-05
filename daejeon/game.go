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

type WinLoose struct {
	Result  string
	Message string
}

var (
	Win   = WinLoose{"Win", "you are winner in this game"}
	Loose = WinLoose{"Loose", "you don't win this game"}
)

type Game struct {
	opts       *Options
	players    map[string]*Player
	joinTicker *time.Ticker
	joinChan   chan *Player
	startChan  chan bool
	finishChan chan []*FinishInfo
	isPlaying  bool
}

func NewGame(opts *Options) *Game {
	g := &Game{
		opts:       opts,
		players:    make(map[string]*Player),
		joinChan:   make(chan *Player, 1),
		startChan:  make(chan bool, 1),
		finishChan: make(chan []*FinishInfo, 1),
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
	counter := g.opts.WaitCountJoin

	for {
		select {
		case p := <-g.joinChan:
			log.Printf("new player \"%s\" joined\n", p.Name)
			g.joinTicker.Reset(time.Second)
			counter = g.opts.WaitCountJoin

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
		log.Println("start a new game")
		g.isPlaying = true

		baseball := NewBaseBall(g.opts.Length, g.opts.Chance)
		log.Println("random number generated", baseball.answer)

		var finishInfos []*FinishInfo
		var hasWinner bool
		for baseball.remainChance > 0 {
			for counter := g.opts.WaitCountGuess; counter > 0; counter-- {
				log.Printf("wait for %d sec to request guessing\n", counter)
				time.Sleep(time.Second)
			}

			baseball.remainChance--
			finishInfos, hasWinner = g.PlayRound(baseball)
			if hasWinner {
				log.Println("finish game because of winner")
				break
			}
		}

		if !hasWinner {
			log.Println("finish game")
		}
		g.finishChan <- finishInfos
	}
}

func (g *Game) finishManager() {
	for {
		finishInfos := <-g.finishChan
		for _, info := range finishInfos {
			go g.sendFinish(info)
		}
		log.Println("the game has been finished")
		g.resetPlayInfo()
	}
}

type FinishInfo struct {
	player     *Player
	answer     string
	usedChance int
	winLoose   WinLoose
}

func (g *Game) resetPlayInfo() {
	g.isPlaying = false
	g.players = make(map[string]*Player)
}

type PlayerGuess struct {
	Name   string
	Number string
}

func (g *Game) PlayRound(baseball *BaseBall) ([]*FinishInfo, bool) {
	players := g.getPlayers()

	playerGuesses := make([]*PlayerGuess, len(players))
	wg := sync.WaitGroup{}
	for i, player := range players {
		wg.Add(1)
		go func(idx int, p *Player) {
			defer wg.Done()
			guessed, err := g.requestGuessing(p, baseball.remainChance)
			if err != nil {
				log.Println(err)
				return
			}
			playerGuesses[idx] = &PlayerGuess{
				Name:   p.Name,
				Number: guessed,
			}
			log.Printf("player: %s, guessed number: %s\n", p.Name, guessed)
		}(i, player)
	}
	wg.Wait()

	resultInfos := make([]*ResultInfo, len(players))
	for i, playerGuess := range playerGuesses {
		if playerGuess == nil {
			log.Printf("player %s failed to guess\n", players[i].Name)
			continue
		}
		result, err := baseball.compareToAnswer(playerGuess.Number)
		if err != nil {
			log.Println(err)
			return nil, false
		}
		resultInfos[i] = &ResultInfo{
			Name:   playerGuess.Name,
			Number: playerGuess.Number,
			Result: result,
		}
		log.Printf("player %s, guessed: %s, result: %+v\n", playerGuess.Name, playerGuess.Number, result)
	}

	for _, player := range players {
		wg.Add(1)
		go func(p *Player) {
			defer wg.Done()
			log.Printf("notifing results to %s\n", p.Name)
			err := g.notifyResults(p, resultInfos)
			if err != nil {
				log.Println(err)
				return
			}
		}(player)
	}
	wg.Wait()

	finInfos := make([]*FinishInfo, len(resultInfos))
	var hasWinner bool
	for i, info := range resultInfos {
		winLoose := Loose
		if info.Result.Win {
			hasWinner = true
			winLoose = Win
		}
		finInfos[i] = &FinishInfo{
			player:     g.players[info.Name],
			answer:     baseball.answer,
			usedChance: baseball.usedChance(),
			winLoose:   winLoose,
		}
	}

	return finInfos, hasWinner
}

type GuessRequest struct {
	Length       int `json:"length"`
	RemainChance int `json:"remain_chance"`
}

type GuessResponse struct {
	Number string `json:"number"`
}

func (g *Game) requestGuessing(p *Player, remainChance int) (string, error) {
	guessReq := GuessRequest{
		Length:       4,
		RemainChance: remainChance,
	}

	b, _ := json.Marshal(guessReq)

	resp, err := sendPost(p.Address.String()+"/guess", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	b, err = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("statusCode: %d, body: %s\n", resp.StatusCode, string(b))
	}

	var guessResp GuessResponse
	err = json.Unmarshal(b, &guessResp)
	if err != nil {
		return "", err
	}
	return guessResp.Number, nil
}

type ResultInfo struct {
	Name   string  `json:"name"`
	Number string  `json:"number"`
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

type FinishRequest struct {
	Name       string `json:"name"`
	Answer     string `json:"answer"`
	UsedChance int    `json:"used_chance"`
	Result     string `json:"result"`
	Message    string `json:"message"`
}

func (g *Game) sendFinish(finInfo *FinishInfo) {
	finReq := FinishRequest{
		Name:       finInfo.player.Name,
		Answer:     finInfo.answer,
		UsedChance: finInfo.usedChance,
		Result:     finInfo.winLoose.Result,
		Message:    finInfo.winLoose.Message,
	}
	b, _ := json.Marshal(finReq)
	log.Println(string(b))

	resp, err := sendPost(finInfo.player.Address.String()+"/finish", bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
}

func (g *Game) getPlayers() []*Player {
	countOfPlayers := len(g.players)
	players := make([]*Player, 0, countOfPlayers)
	for _, player := range g.players {
		players = append(players, player)
	}
	return players
}

func sendPost(url string, data io.Reader) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}
