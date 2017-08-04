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
	"embedit/utils"
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

// Register all of the Services so we can dynamically invoke them
var serviceRegistry = map[string]services.MediaService{
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
		errResp := utils.ErrorMessage{
			Code:    400,
			Message: "No Search Query Provided",
		}
		json.NewEncoder(w).Encode(errResp)
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
		errResp := utils.ErrorMessage{
			Code:    400,
			Message: "No Services Provided",
		}
		json.NewEncoder(w).Encode(errResp)
		return

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
