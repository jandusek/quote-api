package main

import (
	"context"

	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func fetchRandomQuote(client *mongo.Client) []byte {

	quotesCollection := client.Database("quotes").Collection("quotes")
	matchStage := bson.D{{"$sample", bson.D{{"size", 1}}}}
	quoteCursor, err := quotesCollection.Aggregate(context.TODO(), mongo.Pipeline{matchStage})
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = quoteCursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	doc, err := json.Marshal(results[0])
	return doc
}
