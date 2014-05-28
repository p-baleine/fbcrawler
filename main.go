package main

import (
	_ "github.com/lib/pq"

	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	graphUrl string = "https://graph.facebook.com"
)

type paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type feed struct {
	Id string `json:"id"`
}

type graphResponse struct {
	Data []feed `json:"data"`
	Paging paging `json:"paging"`
}

func graph(path string, accessToken string) (resp *http.Response, err error) {
	return http.Get(graphUrl + path + "?access_token=" + accessToken)
}

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")

	if accessToken == "" {
		fmt.Print("access token is not specified")
	}

	resp, err := graph("/10150149727825637/feed", accessToken)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	var data graphResponse
	err = json.Unmarshal(body, &data)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	for _, feed := range data.Data {
		fmt.Printf("%s\n", feed.Id)
	}

	db, err := sql.Open("postgres", "user=postgres dbname=mycrawler")

	defer db.Close()

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	query := "insert into feeds (id) values ($1)"
	stmt, err := db.Prepare(query)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	for _, feed := range data.Data {
		row := stmt.QueryRow(feed.Id)
		fmt.Printf("%s\n", row)
	}

}
