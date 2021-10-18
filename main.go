package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	host       = flag.String("ho", "", "mongo host")
	user       = flag.String("u", "", "mongo user")
	password   = flag.String("p", "", "mongo password")
	db         = flag.String("d", "", "mongo db")
	collection = flag.String("c", "", "mongo collection")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017", *user, *password, *host)))
	if err != nil {
		log.Panicf("failed to connect mongo: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panicf("failed to disconnect mongo: %v", err)
		}
	}()

	collection := client.Database(*db).Collection(*collection)

	filter := bson.D{{"anjie", 1}}

	var result bson.D
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Panicf("failed to search: %v", err)
	}

	fmt.Println(result)
}
