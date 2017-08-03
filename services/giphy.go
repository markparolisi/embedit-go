package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"embedit/utils"
	"embedit/media"
)

type Giphy struct{}

// Make the web request(s) to the Giphy service
// Return a list of the media.Models
func (gi Giphy) GetMedia(q string) ([]media.Model, error) {

	var mediaModels []media.Model

	apiKey, ok := utils.GetConfigValue("giphy", "apiKey")

	// If the required API key for Giphy is not present, return an empty result
	if !ok {
		return mediaModels, fmt.Errorf("Could not get giphy apiKey")
	}

	params := url.Values{}
	params.Add("q", q)
	params.Add("rating", "r")
	params.Add("fmt", "json")
	params.Add("api_key", apiKey)
	reqUrl := fmt.Sprint("https://api.giphy.com/v1/gifs/search?", params.Encode())
	resp, err := http.Get(reqUrl)

	// Return an empty set if the HTTP request fails
	if err != nil {
		return mediaModels, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// Return an empty set if there is a problem with the response body
	if err != nil {
		return mediaModels, err
	}

	// Define Giphy response JSON as a struct
	type giphyResponse struct {
		Data []struct {
			URL     string `json:"url"`
			Slug    string `json:"slug"`
			Created string `json:"import_datetime"`
			Source  string `json:"source"`
			Images map[string]struct {
				URL string `json:"url"`
			} `json:"images"`
		} `json:"data"`
	}

	var gResp giphyResponse

	// Convert the JSOn into the above-defined giphyResponse struct
	err = json.Unmarshal(respBody, &gResp)

	// Return an empty set if the response cannot be unmarshalled
	if err != nil {
		return mediaModels, err
	}

	// Iterate over all of the image structs from the response
	// Convert them into media.Models
	for _, imageData := range gResp.Data {
		created, err := time.Parse("2006-01-02 15:04:05", imageData.Created)

		if err != nil {
			return mediaModels, err
		}

		var thumbnailURL string

		// Get best image in response
		imagePreferences := []string{
			"downsized_small",
			"downsized",
			"original",
		}
		for _, size := range imagePreferences {
			if rendition, ok := imageData.Images[size]; ok && rendition.URL != "" {
				thumbnailURL = rendition.URL
				break
			}
		}
		mediaObject := media.Model{
			Name:         imageData.Slug,
			Service:      "Giphy",
			MediaURL:     imageData.URL,
			Source:       imageData.URL,
			Type:         "gif",
			Created:      created,
			ThumbnailURL: thumbnailURL,
			Credit:       imageData.Source,
		}
		mediaModels = append(mediaModels, mediaObject)
	}

	return mediaModels, nil
}
