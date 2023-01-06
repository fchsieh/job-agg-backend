package database

import (
	"context"
	"log"
	"time"

	"github.com/fchsieh/job-list-backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoInit(c config.Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Mongo.URL))
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
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
		log.Fatalln(err)
		return nil, err
	}

	var results []interface{}
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return results, err
	} else {
		for cur.Next(context.TODO()) {
			var result interface{}
			err := cur.Decode(&result)
			if err != nil {
				log.Fatalln(err)
			}
			results = append(results, result)
		}
	}

	return results, nil
}

func SaveJobsToMongo(conf config.Config, mongo *mongo.Database, date time.Time, jobs []interface{}) error {
	coll := mongo.Collection(conf.Database.Collection)
	_, err := coll.InsertMany(context.TODO(), jobs)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
