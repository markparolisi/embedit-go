package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"embedit/utils"
	"embedit/media"
	"github.com/google/go-querystring/query"
)

type Imgur struct{}

func (im Imgur) getThumbnail(f string) string {

	extension := filepath.Ext(f)
	newFile := f

	if extension != ".gif" {
		newFile = strings.TrimRight(f, extension) + "t" + extension
	}

	return newFile
}

func (im Imgur) GetMedia(q string) ([]media.Model, error) {

	// Each JSON response object
	type dataMedia struct {
		Name     string `json:"title"`
		URL      string `json:"link"`
		Credit   string `json:"account_url"`
		Datetime int64  `json:"datetime"`
	}

	// Hold the JSON response
	type dataResult struct {
		Data []dataMedia `json:"data"`
	}

	type SearchOptions struct {
		Query   string `url:"q"`
		All     string `url:"q_all"`
		Any     string `url:"q_any"`
		Exactly string `url:"q_exactly"`
		Not     string `url:"q_not"`
		Type    string `url:"q_type"`
		SizePx  string `url:"q_size_px"`
	}

	// Initialize return value
	var medias []media.Model

	imgurKey, ok := utils.GetConfigValue("imgur", "clientID")

	if !ok {
		return medias, fmt.Errorf("Could not get imgur clientID")
	}

	mediaTypes := []string{"anigif", "png", "gif"}


	// Using a waitgroup because we have to query each mediaType separately
	wg := sync.WaitGroup{}
	mut := sync.Mutex{}

	for _, mType := range mediaTypes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{}
			response := dataResult{}
			params := SearchOptions{All: q, Type: mType}
			p, err := query.Values(params)
			if err != nil {
				return
			}

			url := fmt.Sprintf("https://api.imgur.com/3/gallery/search/viral?%s", p.Encode())
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return
			}
			req.Header.Add("Authorization", "Client-ID "+imgurKey)
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&response)
			medObjs := make([]media.Model, len(response.Data))

			for i, element := range response.Data {
				httpsURL := strings.Replace(element.URL, "http://", "https://", 1)
				medObjs[i] = media.Model{
					Name:         element.Name,
					Service:      "Imgur",
					Source:       httpsURL,
					Type:         "image",
					Created:      time.Unix(element.Datetime, 0),
					ThumbnailURL: im.getThumbnail(httpsURL),
					MediaURL:     element.URL,
					Credit:       "http://imgur.com/user/" + element.Credit,
				}
			}
			mut.Lock()
			defer mut.Unlock()
			medias = append(medias, medObjs...)
		}()
	}
	wg.Wait()

	return medias, nil

}
