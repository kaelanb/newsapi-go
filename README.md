# News Api Golang
Client for the [News API](https://newsapi.org/)

# Installation
You can simply do `go get github.com/kaelanb/newsapi-go/...`

# DOC
Full Documentation for all the endpoints can be found [here](https://newsapi.org/docs/endpoints)

Request parameters: 

country : ae ar at au be bg br ca ch cn co cu cz de eg fr gb gr 
		  hk hu id ie il in it jp kr lt lv ma mx my ng nl no nz 
		  ph pl pt ro rs ru sa se sg si sk th tr tw ua us ve za 

category: business entertainment general health science sports technology (can't be mixed with sources parameter)

sources: news sources
q: keywords to search for
pageSize: page size
page: page number

# Example
```
import (
	"fmt"
	"github.com/kaelanb/newsapi-go/newsapi"
)

func main() {
	apikey := "apikeyhere"
	n := newsapi.News{}

	query1 := []string{"country=ca"}
	news := n.GetTopHeadlines(apikey, query1...)

	fmt.Println(news)
}
```
