package newsapi

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	testAPIKey = "key"

	topHeadlinesOKResponse = `{
		"status": "ok",
		"totalResults": 1,
		"articles": [{
			"source": {
				"id": "nbc-news",
				"name": "NBC News"
			},
			"author": "Claire Atkinson, Daniella Silva",
			"title": "Giuliani says Trump did not intervene to nix AT&T-Time Warner merger",
			"description": "Rudy Giuliani said Saturday that the president did not intervene in the Department of Justice decision to deny the AT&T- Time Warner merger.",
			"url": "https://www.nbcnews.com/news/us-news/giuliani-suggests-trump-intervened-nix-t-time-warner-merger-n873616",
			"urlToImage": "https://media4.s-nbcnews.com/j/newscms/2018_17/2405161/180419-rudy-giuliani-2016-ac-1140p_14b4f34a294bbd41af891c21ac10aa18.1200;630;7;70;5.jpg",
			"publishedAt": "2018-05-12T19:18:31Z"
		}]
	}`

	errorResponse = `{
		"status": "error",
		"code": "apiKeyMissing",
		"message": "Your API key is missing. Append this to the URL with the apiKey param, or use the x-api-key HTTP header."
	}`
)

func TestClientGetTopHeadlines(t *testing.T) {
	cli := New(testAPIKey)
	publishedAt, _ := time.Parse(time.RFC3339, "2018-05-12T19:18:31Z")

	tests := []struct {
		name     string
		response string
		result   NewsResponse
		wantErr  bool
		err      string
	}{
		{
			name:     "ok response",
			response: topHeadlinesOKResponse,
			result: NewsResponse{
				Status:       "ok",
				TotalResults: 1,
				Articles: []Article{
					Article{
						Source: ArticleSource{
							ID:   "nbc-news",
							Name: "NBC News",
						},
						Author:      "Claire Atkinson, Daniella Silva",
						Title:       "Giuliani says Trump did not intervene to nix AT&T-Time Warner merger",
						Description: "Rudy Giuliani said Saturday that the president did not intervene in the Department of Justice decision to deny the AT&T- Time Warner merger.",
						URL:         "https://www.nbcnews.com/news/us-news/giuliani-suggests-trump-intervened-nix-t-time-warner-merger-n873616",
						URLToImage:  "https://media4.s-nbcnews.com/j/newscms/2018_17/2405161/180419-rudy-giuliani-2016-ac-1140p_14b4f34a294bbd41af891c21ac10aa18.1200;630;7;70;5.jpg",
						PublishedAt: publishedAt,
					},
				},
			},
		},
		{
			name:     "error response",
			response: errorResponse,
			wantErr:  true,
			err:      "code: apiKeyMissing message: Your API key is missing. Append this to the URL with the apiKey param, or use the x-api-key HTTP header.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, tt.response)
			}))

			cli.httpClient.Transport = &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.Dial("tcp", fs.URL[strings.LastIndex(fs.URL, "/")+1:])
				},
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}

			result, err := cli.GetTopHeadlines([]string{"country=ca"})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}

				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("expected error %s, got %s", tt.err, err)
				}

				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			if !reflect.DeepEqual(*result, tt.result) {
				t.Errorf("expected %#v, got %#v", tt.result, result)
			}
		})
	}
}
