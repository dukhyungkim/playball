package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	g := NewGame()
	r := setupRouter(g)

	err := r.Run()
	if err != nil {
		log.Panicln(err)
	}
}

type Message struct {
	Message string `json:"message"`
}

type JoinRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type JoinResponse struct {
	Message string `json:"message"`
}

func setupRouter(g *Game) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, Message{Message: "pong"})
	})

	r.POST("/join", func(c *gin.Context) {
		var joinReq JoinRequest
		err := c.Bind(&joinReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, Message{Message: err.Error()})
			return
		}

		err = g.JoinPlayer(&joinReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, Message{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, JoinResponse{Message: fmt.Sprintf("user %s has been joined the game", joinReq.Name)})
	})

	return r
}
