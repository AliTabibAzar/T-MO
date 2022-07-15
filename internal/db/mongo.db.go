package db

import (
	"context"
	"os"
	"time"

	"github.com/AliTr404/T-MO/pkg/tol"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))

	if err != nil {
		panic("Can't connect to MongoDB!")
	}
	tol.TMessage("MongoDB is starting")
	return client
}
