package server

import (
	"embedit/media"
	"embedit/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Instantiate the route and being the http service
func RunServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/media", mediaIndex)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Just return a simple page directing users to the /media endpoint
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Please access the service at /media")
}

// Interface to collect all of the Services
type mediaService interface {
	GetMedia(string) ([]media.Model, error)
}

// Register all of the Services so we can dymanically invoke them
var serviceRegistry = map[string]mediaService{
	"imgur": services.Imgur{},
	"giphy": services.Giphy{},
}

// Entry point into the main route of the app that serves media.Model responses
func mediaIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// The search string
	query := r.URL.Query().Get("q")

	// Throw an error if request does not contain a search query param
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode("No Search Query Provided"); err != nil {
			panic(err)
		}

		return
	}

	searchQuery := url.QueryEscape(query)
	fmt.Println("Searching: " + searchQuery)

	// Collect all of the requests services from the client here
	var servicesList []string

	// The services (Imgur, Giphy, etc. to query)
	servicesCSV := r.URL.Query().Get("services")

	if servicesCSV != "" {
		servicesList = strings.Split(strings.ToLower(servicesCSV), ",")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode("No Services Provided"); err != nil {
			panic(err)
		}

		return
	}

	fmt.Println("Services selected: " + strings.Join(servicesList, ","))

	// Store the sum of all responses to return as JSON response
	var mediaResponse []media.Model

	// Run all of the Service requests to GetMedia
	wg := sync.WaitGroup{}
	wg.Add(len(servicesList))
	for i := range servicesList {

		// Skip over any unknown services
		if serviceRegistry[servicesList[i]] == nil {
			wg.Done()
		} else {

			go func(i int) {
				fmt.Println("Query service: " + servicesList[i])
				result, err := serviceRegistry[servicesList[i]].GetMedia(searchQuery)
				if err != nil {
					fmt.Println(err)
					return
				}
				wg.Done()
				mediaResponse = append(mediaResponse, result...)
			}(i)
		}
	}

	wg.Wait()

	json.NewEncoder(w).Encode(mediaResponse)

}
