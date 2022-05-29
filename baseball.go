package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func play() {
	const length = 4

	randNum := makeRandomNumber(length)
	fmt.Println("정답:", randNum)
	fmt.Println()

	chance, err := getChance()
	if err != nil {
		log.Println("숫자가 아님")
		return
	}

	for i := 0; i < chance; i++ {
		fmt.Print("숫자 입력: ")
		var num string
		_, err = fmt.Scan(&num)
		if err != nil {
			log.Println("입력이 제대로 되지 않았음")
			return
		}

		_, err = strconv.ParseInt(num, 10, 64)
		if err != nil {
			log.Println("숫자가 아님")
			return
		}

		if len(num) != length {
			log.Printf("%d 자리 숫자 입력해야됨\n", length)
			continue
		}

		var (
			strike = 0
			ball   = 0
		)

		for m := range num {
			for n := range randNum {
				if num[m] == randNum[n] && m == n {
					strike++
					break
				}

				if num[m] == randNum[n] {
					ball++
					break
				}
			}
		}

		if strike == 0 && ball == 0 {
			fmt.Println("out")
			continue
		} else {
			fmt.Printf("strike: %d, ball: %d\n", strike, ball)
		}

		if strike == len(randNum) {
			fmt.Println("you win!")
			break
		}
	}
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

func getChance() (int, error) {
	fmt.Print("기회: ")
	var chance int
	_, err := fmt.Scan(&chance)
	if err != nil {
		return 0, err
	}
	return chance, nil
}
