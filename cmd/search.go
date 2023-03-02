package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

type Response struct {
	Query Query `json:"query"`
}

type Query struct {
	SearchInfo SearchInfo `json:"searchinfo"`
	Search     []Search   `json:"search"`
}

type SearchInfo struct {
	TotalHits int `json:"totalhits"`
}

type Search struct {
	Title  string `json:"title"`
	PageId int    `json:"pageid"`
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search for a source or quote",
	Long:  `search for sources of quotes or quotes by keywords`,
	Run: func(cmd *cobra.Command, args []string) {
		keywords := getKeywords(args)
		response := GetSources(keywords)
		displaySources(response)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getKeywords(args []string) string {
	keywords := ""
	for _, arg := range args {
		keywords += arg + "%20"
	}
	return keywords
}

func GetSources(keywords string) Response {
	url := "https://en.wikiquote.org/w/api.php?action=query&format=json&requestid=&list=search&formatversion=2&srsearch="
	response, err := http.Get(url + keywords)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObj Response
	jsonError := json.Unmarshal(responseData, &responseObj)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	return responseObj
}

func displaySources(data Response) {
	for _, entry := range data.Query.Search {
		fmt.Println(entry.Title)
	}
}
