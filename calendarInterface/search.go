package googleapi

import (
	"context"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

func newSearchService(ctx context.Context, client *http.Client) *customsearch.Service {
	searchSrv, err := customsearch.NewService(ctx, option.WithAPIKey(os.Getenv("CUSTOM_GOOGLE_SEARCH")))
	if err != nil {
		log.Fatalf("Unable to create Custom Search service: %v", err)
	}
	return searchSrv
}

func getImageURL(srv *customsearch.Service, query string) string {
	searchResponse, err := srv.Cse.List().Cx("53efa1d856ad14945").FileType(".jpg").SearchType("image").Q(query).Do()
	if err != nil {
		log.Fatalf("Error executing search: %v", err)
	}

	attachment := searchResponse.Items[1].Link

	return attachment
}