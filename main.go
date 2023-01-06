package main

import (
	"log"

	"github.com/fchsieh/job-list-backend/config"
	"github.com/fchsieh/job-list-backend/database"
	"github.com/fchsieh/job-list-backend/routes"
	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

func read_config() config.Config {
	viper.AddConfigPath("..")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	var conf config.Config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalln(err)
	}
	return conf
}

func set_router(conf config.Config) *gin.Engine {
	fb := database.FirebaseInit(conf)
	mongo := database.MongoInit(conf)
	router := gin.Default()

	router.GET("/jobs", func(c *gin.Context) {
		routes.GetJobs(c, fb, mongo)
	})
	router.GET("/jobs/:date", func(c *gin.Context) {
		routes.GetJobsByDate(c, conf, fb, mongo)
	})

	return router
}

func main() {
	conf := read_config()
	router := set_router(conf)
	router.Run(conf.Server.Host + ":" + conf.Server.Port)
}
