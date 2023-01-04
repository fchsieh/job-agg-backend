package main

import (
	c "github.com/fchsieh/job-list-backend/config"
	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

func read_config() c.Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	var config c.Config
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}

func set_router() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return router
}

func main() {
	config := read_config()
	router := set_router()
	router.Run(config.Server.Host + ":" + config.Server.Port)
}
