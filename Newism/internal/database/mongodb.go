package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Mongo struct {
	Database *mongo.Client
}

const Database = "newism"

func New() *Mongo {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}
	return &Mongo{Database: client}
}

func (m *Mongo) GetCl(name string) *mongo.Collection {
	userCollection := m.Database.Database(Database).Collection(name)
	return userCollection
}

// Connect To database Server
//func Connect() {
//	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://nimaism:pass@nimaism.ucoxrje.mongodb.net/?retryWrites=true&w=majority"))
//	if err != nil {
//		log.Fatalln(err)
//	}
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	err = client.Connect(ctx)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	err = client.Ping(ctx, readpref.Primary())
//	if err != nil {
//		log.Fatalln(err)
//	}
//	Data = client
//}
