package main

import (
	"context"
	"io"
	"net/http"
	"time"
)

type PlayRequest struct {
	Name string `json:"name"`
}

type PlayResponse struct {
	Name         string    `json:"name"`
	StartTime    time.Time `json:"start_time"`
	Length       int       `json:"length"`
	RemainChance int       `json:"remain_chance"`
}

type GuessRequest struct {
	Number string `json:"number"`
}

type GuessNotFinishResponse struct {
	Number string `json:"number"`
	Result struct {
		Out    bool `json:"out"`
		Strike int  `json:"strike"`
		Ball   int  `json:"ball"`
	} `json:"result"`
	RemainChance int `json:"remain_chance"`
}

type GuessFinishResponse struct {
	FinishTime time.Time `json:"finish_time"`
	UsedChance int       `json:"used_chance"`
	Answer     string    `json:"answer"`
	Result     string    `json:"result"`
}

func sendRequest(url string, method string, data io.Reader) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}
