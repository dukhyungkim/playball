package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func sendPlayRequest(host string, player Player) (*PlayResponse, error) {
	const playURI = "/play"

	b, _ := json.Marshal(PlayRequest{Name: player.Name})
	log.Printf("%s -> request body: %s\n", playURI, string(b))

	resp, err := sendRequest(host+playURI, http.MethodPost, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	b, err = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("statusCode: %d, body: %s\n", resp.StatusCode, string(b))
		return nil, ErrUnexpectedStatusCode
	}

	var playResp PlayResponse
	err = json.Unmarshal(b, &playResp)
	if err != nil {
		return nil, err
	}

	log.Printf("%s -> response body:  %s\n", playURI, string(b))
	return &playResp, nil
}

func testPlayRequestAgain(host string, player Player) error {
	const playURI = "/play"

	b, _ := json.Marshal(PlayRequest{Name: player.Name})
	log.Printf("%s -> request body: %s\n", playURI, string(b))

	resp, err := sendRequest(host+playURI, http.MethodPost, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	b, err = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNotAcceptable {
		log.Printf("statusCode: %d, body: %s\n", resp.StatusCode, string(b))
		return ErrUnexpectedStatusCode
	}

	log.Printf("%s -> response body:  %s\n", playURI, string(b))
	return nil
}
