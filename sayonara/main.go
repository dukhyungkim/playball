package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/icrowley/fake"
	"github.com/jessevdk/go-flags"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Options struct {
	Host string `long:"host" description:"server address with port" default:"http://localhost:8000"`
}

type Player struct {
	Name string
}

const (
	WIN   = "승리"
	LOOSE = "패배"
)

func main() {
	opts, err := ParseFlags()
	if err != nil {
		if flags.WroteHelp(err) {
			return
		}
		log.Fatalln(err)
	}

	log.Printf("i will connect server to \"%s\"\n", opts.Host)

	player := Player{Name: fake.FullName()}
	log.Println("let's play new game")
	playResponse, err := sendPlayRequest(opts.Host, player)
	if err != nil {
		log.Fatalln(err)
	}
	err = testPlayRequestAgain(opts.Host, player)
	if err != nil {
		log.Fatalln(err)
	}
	startTime := playResponse.StartTime

	var result string
	var endTime time.Time
	for i := 0; i < playResponse.RemainChance; i++ {
		log.Printf("i will guess your number! ... remain_count: %d\n", playResponse.RemainChance-i)
		time.Sleep(time.Second)
		number := makeRandomNumber(playResponse.Length)
		notFinishResponse, finishResponse, err := sendGuessRequest(opts.Host, number)
		if err != nil {
			log.Fatalln(err)
		}

		if notFinishResponse != nil {
			log.Println("still, i have more chance")
			continue
		}

		if finishResponse != nil {
			log.Println("okay. finish the game.")
			result = finishResponse.Result
			endTime = finishResponse.FinishTime
			break
		}
	}

	switch result {
	case WIN:
		log.Println("yes! i win.")
	case LOOSE:
		log.Println("no... i loose.")
	default:
		log.Printf("result must be %s or %s\n", WIN, LOOSE)
	}

	log.Printf("play time: %s\n", endTime.Sub(startTime).String())
}

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

func ParseFlags() (*Options, error) {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		return nil, err
	}
	return &opts, nil
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

func makeRandomNumber(length int) string {
	var nums = "0123456789"
	var randNumBuilder strings.Builder
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(nums))
		randNumBuilder.WriteRune(rune(nums[idx]))
		nums = nums[:idx] + nums[idx+1:]
	}
	randNum := randNumBuilder.String()
	return randNum
}
