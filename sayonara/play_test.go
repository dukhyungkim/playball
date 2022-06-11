package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func Test_sendPlayRequest(t *testing.T) {
	name := fake.LastName()
	want := &PlayResponse{
		Name:         name,
		StartTime:    time.Now().Truncate(time.Nanosecond),
		Length:       rand.Intn(10),
		RemainChance: 5,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, _ := json.Marshal(want)
		if _, err := w.Write(b); err != nil {
			log.Println(err)
			return
		}
	}))
	defer server.Close()

	type args struct {
		host   string
		player Player
	}
	tests := []struct {
		name    string
		args    args
		want    *PlayResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				host: server.URL,
				player: Player{
					Name: name,
				},
			},
			want:    want,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sendPlayRequest(tt.args.host, tt.args.player)
			if !tt.wantErr(t, err, fmt.Sprintf("sendPlayRequest(%v, %v)", tt.args.host, tt.args.player)) {
				return
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("sendPlayRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testPlayRequestAgain(t *testing.T) {
	type args struct {
		host   string
		player Player
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, testPlayRequestAgain(tt.args.host, tt.args.player), fmt.Sprintf("testPlayRequestAgain(%v, %v)", tt.args.host, tt.args.player))
		})
	}
}
