package newsapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const (
	baseURL        = "https://newsapi.org/v2"
	defaultTimeout = 2 * time.Second
)

// Error json struct for error message.
type Error struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements error interface.
func (err Error) Error() string {
	return fmt.Sprintf("code: %s message: %s", err.Code, err.Nessage)
}

// NewsResponse json struct for news articles.
type NewsResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

// Article struct for unmarshaling article
// json data.
type Article struct {
	Source      ArticleSource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	URLToImage  string        `json:"urlToImage"`
	PublishedAt time.Time     `json:"publishedAt"`
}

// ArticleSource struct for unmarshaling source
// data embeded in article json.
type ArticleSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SourcesResponse json struct for source.
type SourcesResponse struct {
	Status  string   `json:"status"`
	Sources []Source `json:"sources"`
}

// Source struct for unmarshaling source
// json data.
type Source struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	Country     string `json:"country"`
}

// Client for api requests.
type Client struct {
	httpClient *http.Client
}

// New returns initialized client.
func New(apiKey string) *Client {
	cli := Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			apiKey:  apiKey,
		},
	}

	return &cli
}

// GetTopHeadlines get top news headlines relevant to given parameters.
func (cli *Client) GetTopHeadlines(args ...string) (*NewsResponse, error) {
	url := buildURL(fmt.Sprintf("%s/%s?", baseURL, "top-headlines"), args...)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", cli.apiKey)

	res, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client get request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read request body")
	}
	defer res.Body.Close()

	nr := NewsResponse{}
	if err := json.Unmarshal(body, &nr); err != nil {
		return nil, errors.Wrap(err, "unmrashal to news response")
	}
	if nr.Status == "error" {
		err := Error{}
		if err := json.Unmarshal(body, &err); err != nil {
			return nil, errors.Wrap(err, "unmarshal to error")
		}
		return nil, err
	}

	return nr, nil

}

// GetEverything Get every article relevant to given parameters.
func (cli *Client) GetEverything(args ...string) (*NewsResponse, error) {
	url := buildURL(fmt.Sprintf("%s/%s?", baseURL, "everything"), args...)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", cli.apiKey)

	res, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client get request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read request body")
	}
	defer res.Body.Close()

	nr := NewsResponse{}
	if err := json.Unmarshal(body, &nr); err != nil {
		return nil, errors.Wrap(err, "unmrashal to news response")
	}
	if nr.Status == "error" {
		err := Error{}
		if err := json.Unmarshal(body, &err); err != nil {
			return nil, errors.Wrap(err, "unmarshal to error")
		}
		return nil, err
	}

	return nr, nil
}

// GetSources get all sources relevant to given parameters.
func (cli *Cliebt) GetSources(args ...string) (*SourcesResponse, error) {
	url := buildURL(fmt.Sprintf("%s/%s?", baseURL, "sources"), args...)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", cli.apiKey)

	res, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client get request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read request body")
	}
	defer res.Body.Close()

	sr := SourcesResponse{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, errors.Wrap(err, "unmrashal to sources response")
	}
	if sr.Status == "error" {
		err := Error{}
		if err := json.Unmarshal(body, &err); err != nil {
			return nil, errors.Wrap(err, "unmarshal to error")
		}
		return nil, err
	}

	return nr, nil

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
