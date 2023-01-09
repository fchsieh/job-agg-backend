package main

import (
	"log"
	"time"

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

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func set_router(conf config.Config) *gin.Engine {
	fb := database.FirebaseInit(conf)
	mongo := database.MongoInit(conf)
	router := gin.Default()
	router.Use(CORS())

	router.GET("/jobs", func(c *gin.Context) {
		routes.GetJobs(c, fb, mongo)
	})
	router.GET("/jobs/:date", func(c *gin.Context) {
		routes.GetJobsByDate(c, conf, fb, mongo)
	})

	return router
}

func del_old_data(conf config.Config) {
	if !conf.Server.DeleteOldData {
		return
	}
	// delete data older than * days
	to_del := time.Now().AddDate(0, 0, -conf.Server.DeleteOldDataAfterDays)

	for {
		time.Sleep(24 * time.Hour)
		err := database.DeleteFirebaseOldData(conf, database.FirebaseInit(conf), to_del)
		if err != nil {
			panic(err)
		} else {
			log.Println("Deleted old data")
		}
		err = database.DeleteMongoOldData(conf, database.MongoInit(conf), to_del)
		if err != nil {
			panic(err)
		} else {
			log.Println("Deleted old data")
		}
	}
}

func main() {
	conf := read_config()
	go del_old_data(conf)
	router := set_router(conf)
	router.Run(conf.Server.Host + ":" + conf.Server.Port)
}
