package server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"net/url"
	"strings"
	"sync"
	"embedit/media"
	"embedit/services"
)

func RunServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/media", mediaIndex)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Please access the service at /media")
}

type mediaService interface {
	GetMedia(string) ([]media.Model, error)
}

var serviceRegistry = map[string]mediaService{
	"imgur": services.Imgur{},
	"giphy": services.Giphy{},
}

func mediaIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// The search string
	query, ok := r.URL.Query()["q"]

	// Throw an error if request does not contain a search query param
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode("No Search Query Provided"); err != nil {
			panic(err)
		}

		return
	}

	searchQuery := url.QueryEscape(query[0])
	fmt.Println("Searching: " + searchQuery)

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

	wg := sync.WaitGroup{}
	wg.Add(len(servicesList))
	for i := range servicesList {

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
