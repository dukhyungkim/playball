package main

import "github.com/jessevdk/go-flags"

type Options struct {
	Host string `long:"host" description:"server address with port" default:"http://localhost:8080"`
}

func ParseFlags() (*Options, error) {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		return nil, err
	}
	return &opts, nil
}
