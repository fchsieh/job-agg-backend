package main

import (
	"net/http"

	c "./config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func read_config() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	var config c.Config
}

func main() {
	read_config()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run("192.168.0.22:8888")
}
