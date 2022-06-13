package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sendGuessRequest(t *testing.T) {
	number := makeRandomNumber(rand.Intn(10))
	want := &GuessNotFinishResponse{
		Number: number,
		Result: struct {
			Out    bool `json:"out"`
			Strike int  `json:"strike"`
			Ball   int  `json:"ball"`
		}{
			Out:    false,
			Strike: 1,
			Ball:   2,
		},
		RemainChance: 5,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		all, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			if err = req.Body.Close(); err != nil {
				log.Println(err)
			}
		}()

		var request GuessRequest
		err = json.Unmarshal(all, &request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, _ := json.Marshal(want)
		if _, err = w.Write(b); err != nil {
			log.Println(err)
			return
		}
	}))

	type args struct {
		host   string
		number string
	}
	tests := []struct {
		name    string
		args    args
		want    *GuessNotFinishResponse
		want1   *GuessFinishResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success but not end",
			args: args{
				host:   server.URL,
				number: number,
			},
			want:    want,
			want1:   nil,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := sendGuessRequest(tt.args.host, tt.args.number)
			if !tt.wantErr(t, err, fmt.Sprintf("sendGuessRequest(%v, %v)", tt.args.host, tt.args.number)) {
				return
			}
			assert.Equalf(t, tt.want, got, "sendGuessRequest(%v, %v)", tt.args.host, tt.args.number)
			assert.Equalf(t, tt.want1, got1, "sendGuessRequest(%v, %v)", tt.args.host, tt.args.number)
		})
	}
}
