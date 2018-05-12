package newsapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//Error json struct for error message
type Error struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

//News json struct for news articles
type News struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		URLToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
	} `json:"articles"`
}

//Source json struct for source
type Source struct {
	Status  string `json:"status"`
	Sources []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Category    string `json:"category"`
		Language    string `json:"language"`
		Country     string `json:"country"`
	} `json:"sources"`
}

//GetTopHeadlines get top news headlines relevant to given parameters
func (n *News) GetTopHeadlines(apikey string, args ...string) *News {
	base := "https://newsapi.org/v2/top-headlines?"

	newsClient := &http.Client{
		Timeout: time.Second * 2,
	}

	url := buildURL(base, args...)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", apikey)

	res, err := newsClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	json.Unmarshal(body, &n)
	if n.Status == "error" {
		jsonErr := Error{}
		json.Unmarshal(body, &jsonErr)
		log.Fatal(jsonErr)
	}

	return n
}

//GetEverything Get every article relevant to given parameters
func (n *News) GetEverything(apikey string, args ...string) *News {
	base := "https://newsapi.org/v2/everything?"

	newsClient := &http.Client{
		Timeout: time.Second * 2,
	}

	url := buildURL(base, args...)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", apikey)

	res, err := newsClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	json.Unmarshal(body, &n)
	if n.Status == "error" {
		jsonErr := Error{}
		json.Unmarshal(body, &jsonErr)
		log.Fatal(jsonErr)
	}

	return n
}

//GetSources get all sources relevant to given parameters
func (s *Source) GetSources(apikey string, args ...string) *Source {
	base := "https://newsapi.org/v2/sources?"

	newsClient := &http.Client{
		Timeout: time.Second * 2,
	}

	url := buildURL(base, args...)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", apikey)

	res, err := newsClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	json.Unmarshal(body, &s)
	if s.Status == "error" {
		jsonErr := Error{}
		json.Unmarshal(body, &jsonErr)
		log.Fatal(jsonErr)
	}

	return s
}

func buildURL(base string, args ...string) (escapedURL string) {
	var params string
	for i, param := range args {
		params += param
		if i < len(args)-1 {
			params += "&"
		}
	}
	p := url.PathEscape(params)
	u := base + p

	return u
}
