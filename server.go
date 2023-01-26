package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	anyascii "github.com/anyascii/go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get_quote" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	quote := fetchRandomQuote(client)
	fmt.Fprintf(w, "%s\n", quote)
}

func getQuoteHandlerAscii(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get_quote_ascii" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var quoteJson []byte = fetchRandomQuote(client)

	// decode JSON first
	var quote map[string]any
	json.Unmarshal([]byte(quoteJson), &quote)

	// then transliterate
	type Quote struct {
		Id       string `json:"_id"`
		Quote    string `json:"quote"`
		Followup string `json:"followup,omitempty"`
		Author   string `json:"author"`
	}
	newQuote := &Quote{
		Author: quote["author"].(string),
		Quote:  anyascii.Transliterate(quote["quote"].(string)),
		Id:     quote["_id"].(string),
	}
	if quote["followup"] != nil {
		newQuote.Followup = anyascii.Transliterate(quote["followup"].(string))
	}

	// then encode JSON again
	data, _ := json.Marshal(newQuote)
	fmt.Fprintf(w, "%s\n", string(data))
}

func main() {
	dbConnString := os.Getenv("DB_CONN_STRING")
	if dbConnString == "" {
		log.Fatal("Missing DB_CONN_STRING")
	}
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dbConnString))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	http.HandleFunc("/get_quote", getQuoteHandler)
	http.HandleFunc("/get_quote_ascii", getQuoteHandlerAscii)

	fmt.Printf("Starting server at port %s\n", httpPort)
	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		log.Fatal(err)
	}
}
