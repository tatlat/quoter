package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/spf13/cobra"
)

type Page struct {
	Parse Parse `json:"parse"`
}

type Parse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Source struct {
	title  string
	quotes []string
}

var (
	max_quotes int
	random     bool
)

// fromCmd represents the from command
var fromCmd = &cobra.Command{
	Use:   "from",
	Short: "get a quote from a page that matches the args",
	Long:  `obtains the first quote from a page that closest matches the supplied arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		keywords := getKeywords(args)
		response := GetSources(keywords)
		pageId := response.Query.Search[0].PageId
		page := getPage(pageId)
		source := getSource(page.Parse.Text, page.Parse.Title)
		displayQuotes(source)
	},
}

func init() {
	rootCmd.AddCommand(fromCmd)

	fromCmd.Flags().IntVarP(&max_quotes, "max_quotes", "mq", 1, "The number of quotes to retrieve")
	fromCmd.Flags().BoolVarP(&random, "random", "r", false, "Randomizes quote selection if true")
}

func getPage(pageId int) Page {
	url := "https://en.wikiquote.org/w/api.php?action=parse&format=json&requestid=&formatversion=2&pageid="
	fullUrl := fmt.Sprint(url, pageId)
	response, err := http.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObj Page
	jsonError := json.Unmarshal(responseData, &responseObj)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	return responseObj
}

func getSource(text string, title string) Source {
	var source Source
	source.title = title
	source.quotes = getQuotes(text)
	return source
}

func getQuotes(text string) []string {
	quotes := make([]string, 0)
	tokenizer := html.NewTokenizer(strings.NewReader(text))

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}

			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		}

		if tokenType == html.StartTagToken {
			//get the token
			token := tokenizer.Token()
			if "ul" == token.Data {
				li := extractList(tokenizer)
				if len(li) == 0 {
					continue
				}
				quotes = append(quotes, li)
			} else if "dd" == token.Data {
				q := extractDD(tokenizer)
				if len(q) == 0 || strings.Contains(q, "Seasons") || !strings.Contains(q, ":") || strings.Contains(q, "See also:") {
					continue
				}
				quotes = append(quotes, q)
			}
		}
	}

	return quotes
}

func extractList(tokenizer *html.Tokenizer) string {
	tokenType := tokenizer.Next()
	retVal := ""

	if tokenType == html.StartTagToken {
		token := tokenizer.Token()
		if "li" == token.Data {

			for {
				tokenType = tokenizer.Next()

				if tokenType == html.EndTagToken && tokenizer.Token().Data == "ul" {
					break
				}

				if tokenType == html.TextToken {
					retVal += tokenizer.Token().Data
				}
			}
		}
	}
	return retVal
}

func extractDD(tokenizer *html.Tokenizer) string {
	tokenType := tokenizer.Next()
	retVal := ""

	for {
		if tokenType == html.EndTagToken && tokenizer.Token().Data == "dd" {
			break
		}

		if tokenType == html.TextToken {
			text := tokenizer.Token().Data
			retVal += text
		}

		tokenType = tokenizer.Next()
	}

	return retVal
}

func displayQuotes(source Source) {
	fmt.Println("Source: " + source.title)
	fmt.Println(source.quotes[0])
}
