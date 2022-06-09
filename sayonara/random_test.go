package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeRandomNumber(t *testing.T) {
	const length = 4
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "must be a unique number",
			args: args{
				length: length,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeRandomNumber(tt.args.length)
			assert.Equal(t, len(got), length)
			assert.Equal(t, hasDuplicatedNumber(got), false)
		})
	}
}

func hasDuplicatedNumber(num string) bool {
	split := strings.Split(num, "")

	m := make(map[string]struct{})
	for _, s := range split {
		m[s] = struct{}{}
	}

	return len(num) != len(m)
}
