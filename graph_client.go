package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Graph api client.
type GraphClient struct {
	graphUrl string
	accessToken string
}

// Paging field of graph api response.
type paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

// Feed field of graph api response.
type feed struct {
	Id string `json:"id"`
}

// Graph Api response.
type GraphResponse struct {
	Data []feed `json:"data"`
	Paging paging `json:"paging"`
}

// Request `path` via GET method.
// Returns Graph Api response.
func (client *GraphClient) Get(path string) (*GraphResponse, error) {
	params := url.Values{}
	params.Add("access_token", client.accessToken)

	resp, err := http.Get(client.graphUrl + path + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data *GraphResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
