package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/fchsieh/job-list-backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func IsRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

func MongoInit(c config.Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Select the correct mongo URI based on whether we are running in docker or not
	var mongoURI string
	if IsRunningInDocker() {
		mongoURI = c.Mongo.DockerURI
	} else {
		mongoURI = c.Mongo.LocalURI
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Println("Failed to connect to mongo")
		log.Fatalln(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("Failed to ping mongo")
		log.Fatalln(err)
	}
	return client.Database(c.Database.Database)
}

func FetchMongoJobsByDate(conf config.Config, mongo *mongo.Database, date time.Time) ([]interface{}, error) {
	// Construct a query to fetch the jobs for the given date
	coll := mongo.Collection(conf.Database.Collection)
	cur, err := coll.Find(context.TODO(), map[string]interface{}{
		"date_posted": date.Format("2006-01-02"),
	})
	if err != nil {
		return nil, err
	}

	var results []interface{}
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return results, err
	} else {
		for cur.Next(context.TODO()) {
			var result config.Job
			err := cur.Decode(&result)
			if err != nil { // failed to decode from mongo
				log.Fatalln(err)
			}
			results = append(results, result)
		}
	}
	return results, nil
}

func SaveJobsToMongo(conf config.Config, mongo *mongo.Database, date time.Time, jobs []interface{}) error {
	coll := mongo.Collection(conf.Database.Collection)
	var jobsToSave []interface{}
	for _, job := range jobs {
		job := job.(map[string]interface{})
		job["date_posted"] = date.Format("2006-01-02")
		jobsToSave = append(jobsToSave, job)
	}
	_, err := coll.InsertMany(context.TODO(), jobsToSave)
	if err != nil { // failed to insert to mongo
		log.Fatalln(err)
		return err
	}

	return nil
}

func DeleteMongoOldData(conf config.Config, mongo *mongo.Database, date time.Time) error {
	coll := mongo.Collection(conf.Database.Collection)
	_, err := coll.DeleteMany(context.TODO(), map[string]interface{}{
		"date_posted": date.Format("2006-01-02"),
	})
	if err != nil {
		return err
	}
	return nil
}
