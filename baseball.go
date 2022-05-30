package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type BaseBall struct {
	answer string
	length int
}

func NewBaseBall(length int) *BaseBall {
	return &BaseBall{
		answer: makeRandomNumber(length),
		length: length,
	}
}

type Result struct {
	Win    bool `json:"win"`
	Out    bool `json:"out"`
	Strike int  `json:"strike"`
	Ball   int  `json:"ball"`
}

func (b *BaseBall) compareToAnswer(number int) (*Result, error) {
	var (
		strike = 0
		ball   = 0
	)

	guessed := fmt.Sprint(number)
	if len(guessed) != b.length {
		return nil, ErrLengthMismatched
	}

	if guessed == b.answer {
		return &Result{Win: true}, nil
	}

	for m := range guessed {
		for n := range b.answer {
			if guessed[m] == b.answer[n] && m == n {
				strike++
				break
			}

			if guessed[m] == b.answer[n] {
				ball++
				break
			}
		}
	}

	result := Result{}
	if strike == 0 && ball == 0 {
		result.Out = true
	} else {
		result.Strike = strike
		result.Ball = ball
	}

	return &result, nil
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
