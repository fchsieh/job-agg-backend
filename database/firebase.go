package database

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/fchsieh/job-list-backend/config"
	"google.golang.org/api/option"
)

func FirebaseInit(c config.Config) *firestore.Client {
	// Set the configuration file
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: c.Firebase.URL,
	}
	if _, err := os.Stat(c.Firebase.Credentials); os.IsNotExist(err) {
		if _, err := os.Stat("../" + c.Firebase.Credentials); os.IsNotExist(err) {
			log.Fatalln(err)
		}
		// copy the file to the current directory
		if err := os.Link("../"+c.Firebase.Credentials, c.Firebase.Credentials); err != nil {
			log.Fatalln(err)
		}
	}
	opt := option.WithCredentialsFile(c.Firebase.Credentials)
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return client
}

func FetchFirebaseJobsByDate(conf config.Config, fb *firestore.Client, date time.Time) ([]interface{}, error) {
	// Construct a query to fetch the jobs for the given date
	data, err := fb.Collection(conf.Database.Collection).Doc(date.Format("2006-01-02")).Get(context.Background())
	if err != nil {
		return nil, err
	}
	jobs := data.Data()[conf.Database.Document]
	// Convert the jobs to a slice of Job structs
	fbJobs := jobs.([]interface{})
	return fbJobs, nil
}

func DeleteFirebaseOldData(conf config.Config, fb *firestore.Client, date time.Time) error {
	// Construct a query to fetch the jobs for the given date
	_, err := fb.Collection(conf.Database.Collection).Doc(date.Format("2006-01-02")).Delete(context.Background())
	if err != nil {
		return err
	}
	return nil
}
