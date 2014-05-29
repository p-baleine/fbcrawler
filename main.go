package main

import (
	_ "github.com/lib/pq"

	"database/sql"
	"fmt"
	"os"
)

// Graph api url.
const graphUrl string = "https://graph.facebook.com"

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")

	if accessToken == "" {
		panic(fmt.Sprint("access token is not specified"))
	}

	client := &GraphClient{graphUrl, accessToken}
	data, err := client.Get("/10150149727825637/feed")

	if err != nil {
		panic(fmt.Sprintf("%s\n", err))
	}

	for _, feed := range data.Data {
		fmt.Printf("%s\n", feed.Id)
	}

	// get db connection
	db, err := sql.Open("postgres", "user=postgres dbname=mycrawler")
	if err != nil {
		panic(fmt.Sprintf("%s\n", err))
	}
	defer db.Close()

	// bulk insert
	query := "insert into feeds (id) values ($1)"
	stmt, err := db.Prepare(query)

	if err != nil {
		panic(fmt.Sprintf("%s\n", err))
	}

	for _, feed := range data.Data {
		row := stmt.QueryRow(feed.Id)
		fmt.Printf("%s\n", row)
	}

}
