package main

import (
	"log"
	"net/url"
)

type Player struct {
	Name    string
	Address *url.URL
}

func NewPlayer(req *JoinRequest) (*Player, error) {
	if req.Name == "" {
		return nil, ErrEmptyName
	}

	if req.Address == "" {
		return nil, ErrEmptyAddress
	}

	parsedURL, err := url.Parse(req.Address)
	if err != nil {
		log.Println(err)
		return nil, ErrInvalidAddress
	}

	if parsedURL.Scheme == "" {
		return nil, ErrInvalidAddress
	}

	return &Player{
		Name:    req.Name,
		Address: parsedURL,
	}, nil
}
