package main

import (
	"fmt"
	"os"

	"github.com/wolv89/troster/internal/cache"
	"github.com/wolv89/troster/internal/cli"
	"github.com/wolv89/troster/internal/models"
	"github.com/wolv89/troster/internal/scraper"
	"github.com/wolv89/troster/internal/server"
)

func main() {
	flags := cli.ParseFlags()
	cfg := cli.GetConfig(flags)

	comps, err := scraper.FetchCompetitions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching competitions: %v\n", err)
		os.Exit(1)
	}

	if len(comps) == 0 {
		fmt.Fprintln(os.Stderr, "No competitions found on TROLS")
		os.Exit(1)
	}

	comp := cli.ChooseCompetition(comps)

	var data *models.ScrapedData
	if !flags.ForceScrape {
		cached, err := cache.Load(comp.Value)
		if err == nil {
			fmt.Printf("Using cached data (scraped at %s)\n\n", cached.ScrapedAt)
			data = cached
		}
	}

	if data == nil {
		scraped, err := scraper.ScrapeAll(comp, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scraping fixtures: %v\n", err)
			os.Exit(1)
		}
		data = scraped

		if err := cache.Save(comp.Value, data); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not cache data: %v\n", err)
		}
	}

	if err := server.Start(data, 8080); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
