package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func sendGuessRequest(host, number string) (*GuessNotFinishResponse, *GuessFinishResponse, error) {
	const guessURI = "/guess"

	b, _ := json.Marshal(GuessRequest{Number: number})
	log.Printf("%s -> request body: %s\n", guessURI, string(b))

	resp, err := sendRequest(host+guessURI, http.MethodPut, bytes.NewBuffer(b))
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	b, err = io.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusOK:
		var guessResp GuessNotFinishResponse
		err = json.Unmarshal(b, &guessResp)
		if err != nil {
			return nil, nil, err
		}

		log.Printf("%s -> response body:  %s\n", guessURI, string(b))
		return &guessResp, nil, nil

	case http.StatusCreated:
		var guessResp GuessFinishResponse
		err = json.Unmarshal(b, &guessResp)
		if err != nil {
			return nil, nil, err
		}

		log.Printf("%s -> response body:  %s\n", guessURI, string(b))
		return nil, &guessResp, nil

	default:
		log.Printf("statusCode: %d, body: %s\n", resp.StatusCode, string(b))
		return nil, nil, errors.New("unexpected status code")
	}
}
