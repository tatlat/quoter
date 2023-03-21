package cmd

import (
	"quoter/tui"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "displays a dashboard of quotes",
	Long:  `displays a dashboard of quotes from favorite quotes defined in config file`,
	Run: func(cmd *cobra.Command, args []string) {
		favorites := loadFavorites()
		quotes := getFaveQuotes(favorites)
		tui.LoadUI(quotes)
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

func loadFavorites() []string {
	favorites := viper.GetStringSlice("favorites")
	return favorites[:]
}

func getFaveQuotes(favorites []string) []string {
	quotes := make([]string, 0)
	for _, fave := range favorites {
		fave := strings.ReplaceAll(fave, " ", "%20")
		sources := GetSources(fave)
		pageId := sources.Query.Search[0].PageId
		page := getPage(pageId)
		source := getSource(page.Parse.Text, page.Parse.Title)
		quotes = append(quotes, source.quotes[0])
	}
	return quotes
}
