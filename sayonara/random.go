package main

import (
	"math/rand"
	"strings"
)

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
