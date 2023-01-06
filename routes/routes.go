package routes

import (
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/fchsieh/job-list-backend/config"
	"github.com/fchsieh/job-list-backend/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetJobs(c *gin.Context, fb *firestore.Client, mongo *mongo.Database) {
	c.JSON(200, gin.H{
		"message": "jobs",
		"jobs":    "jobs",
	})
}

func GetJobsByDate(c *gin.Context, conf config.Config, fb *firestore.Client, mongo *mongo.Database) {
	dateStr := c.Param("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
			"error":   "date is empty",
		})
		return
	}
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'date' query parameter"})
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'date' query parameter"})
		return
	}
	// Try to find jobs by date in mongo
	mongoJobs, err := database.FetchMongoJobsByDate(conf, mongo, date)
	if err != nil || len(mongoJobs) == 0 {
		fbJobs, err := database.FetchFirebaseJobsByDate(conf, fb, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching jobs from firebase"})
			return
		}
		// save jobs to mongo
		err = database.SaveJobsToMongo(conf, mongo, date, fbJobs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving jobs to mongo"})
			return
		}
		// done, return jobs
		c.JSON(http.StatusOK, gin.H{"jobs": fbJobs})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": mongoJobs})
}
