package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/icrowley/fake"
	"github.com/jessevdk/go-flags"
)

func init() {
	rand.Seed(time.Now().Unix())
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
	log.Println("let's play a new game")
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
