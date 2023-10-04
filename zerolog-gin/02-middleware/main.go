package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {

	r := gin.New()
	r.Use(GinLogger(log.Logger), GinRecovery(log.Logger, true))

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Hello": "Lixin",
		})
	})

	r.GET("/world", func(c *gin.Context) {
		panic(errors.New("Panic /world"))
	})

	r.Run()
}
