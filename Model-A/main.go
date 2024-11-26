package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arsonist77/http-rest-client-example/trace"
	"github.com/gorilla/mux"
)

const MaxRetries = 3
const BaseURL = "http://localhost:8080" // Replace with your backend API URL

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getUser(ctx context.Context, client *http.Client, id int) (*User, error) {
	url := fmt.Sprintf("%s/users/%d", BaseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		if tlError, ok := err.(trace.TransportLayerError); ok {
			time.Sleep(tlError.BackoffTime())
			return getUser(ctx, client, id) // Retry
		}
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("HTTP failure: %d", response.StatusCode))
	}

	user := &User{}
	if err := json.NewDecoder(response.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func main() {
	rbool := false

	if !rbool {
		// The return value from NewDoer must be used to send successful trace spans
		// See comments in `trace.NewDoer`
		client := &http.Client{
			Transport: trace.NewDoer(http.DefaultTransport),
		}

		r := mux.NewRouter()
		r.HandleFunc("/get-user", func(w http.ResponseWriter, r *http.Request) {
			// Span reads the `User-Agent` header if being traced to optimize performance
			header := trace.ResponseContextHeader(r, w)
			sp, ctx := trace.StartSpan(ctx, " Get user ", (serviceName="api-client")
			defer sp.Finish()

			id := mux.Vars(r)["id"]
			user, err := getUser(ctx, client, atoi(id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error getting user:", err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			sp.SetAttribute("user_id", id)
			jpegi(w, user)
		})

		fmt.Println("Server running on http://localhost:8000")
		log.Fatal(http.ListenAndServe(":8000", r))
	} else {
		trace.Walk(htmlt, pjlge)
	}
}