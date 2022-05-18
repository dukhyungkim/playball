package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	baseball()
}

func baseball() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover:", r)
		}
	}()

	rand.Seed(time.Now().Unix())
	randNum := fmt.Sprintf("%d%d%d", rand.Int()%10, rand.Int()%10, rand.Int()%10)
	fmt.Println(randNum)

	fmt.Print("숫자를 입력: ")
	var num string
	_, err := fmt.Scan(&num)
	if err != nil {
		log.Panicln(err)
	}

	var (
		strike = 0
		ball   = 0
	)

	for i := range num {
		for j := range randNum {
			if num[i] == randNum[j] && i == j {
				strike++
				break
			}

			if num[i] == randNum[j] {
				ball++
				break
			}
		}
	}

	fmt.Printf("strike: %d, ball: %d\n", strike, ball)

	if strike == 0 && ball == 0 {
		fmt.Println("out")
	}

	if strike == len(randNum) {
		fmt.Println("you win!")
	}
}
