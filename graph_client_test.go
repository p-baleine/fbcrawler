package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGraphClientGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		token := query["access_token"][0]
		if  token != "dummyaccesstoken" {
			t.Errorf("query parameter access_token = %q, want %q", token, "dummyaccesstoken")
		}

		fmt.Fprintln(w, `
{
  "data": [
    { "id": "123_456"}
  ],
  "paging": {
    "previous": "previous",
    "next": "next"
  }
}
    `)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	client := &GraphClient{ts.URL, "dummyaccesstoken"}
	resp, err := client.Get("/1234/feed")

	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("Neither `err` and `resp` should be nil.")
	}

	if resp.Data[0].Id != "123_456" {
		t.Errorf("Get() = %q, want %q", "123_456", resp.Data[0].Id)
	}
}
