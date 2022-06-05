package main

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Length         int `long:"length" description:"length of answer number" default:"4"`
	Chance         int `long:"chance" description:"chance of trying to guess" default:"5"`
	WaitCountJoin  int `long:"join_wait" description:"count of waiting count to join a new player" default:"10"`
	WaitCountGuess int `long:"guess_wait" description:"count of waiting count before guessing" default:"10"`
}

func ParseFlags() (*Options, error) {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		return nil, err
	}
	return &opts, nil
}
