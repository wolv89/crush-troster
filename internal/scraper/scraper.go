package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wolv89/troster/internal/models"
	"github.com/wolv89/troster/internal/parser"
)

const (
	baseURL  = "https://trols.org.au/brta/fixture.php"
	delayMs  = 500
)

func post(params url.Values) (string, error) {
	resp, err := http.Post(baseURL, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return "", fmt.Errorf("POST failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	return string(body), nil
}

func FetchCompetitions() ([]models.Competition, error) {
	html, err := post(url.Values{})
	if err != nil {
		return nil, err
	}
	return parser.ParseCompetitions(html), nil
}

func FetchSections(compValue string) ([]models.Section, error) {
	time.Sleep(delayMs * time.Millisecond)

	html, err := post(url.Values{
		"daytime": {compValue},
		"which":   {"0"},
		"style":   {""},
	})
	if err != nil {
		return nil, err
	}
	return parser.ParseSections(html), nil
}

func FetchTeams(compValue, sectionValue string) ([]models.Team, error) {
	time.Sleep(delayMs * time.Millisecond)

	html, err := post(url.Values{
		"daytime": {compValue},
		"section": {sectionValue},
		"which":   {"1"},
		"style":   {""},
	})
	if err != nil {
		return nil, err
	}

	teams := parser.ParseTeams(html, sectionValue)
	return teams, nil
}

func FetchFixture(compValue, sectionValue, teamValue string) ([]models.FixtureRound, error) {
	time.Sleep(delayMs * time.Millisecond)

	html, err := post(url.Values{
		"daytime": {compValue},
		"section": {sectionValue},
		"team":    {teamValue},
		"which":   {"2"},
		"style":   {""},
	})
	if err != nil {
		return nil, err
	}

	return parser.ParseFixture(html, teamValue), nil
}

func ScrapeAll(comp models.Competition, cfg models.Config) (*models.ScrapedData, error) {
	fmt.Printf("Scraping competition: %s\n", comp.Label)

	sections, err := FetchSections(comp.Value)
	if err != nil {
		return nil, fmt.Errorf("fetching sections: %w", err)
	}
	fmt.Printf("Found %d sections\n", len(sections))

	data := &models.ScrapedData{
		Competition: comp,
		Config:      cfg,
		ScrapedAt:   time.Now().Format(time.RFC3339),
	}

	for _, section := range sections {
		fmt.Printf("  Section: %s\n", section.Label)

		teams, err := FetchTeams(comp.Value, section.Value)
		if err != nil {
			fmt.Printf("    Error fetching teams: %v\n", err)
			continue
		}

		clubTeams := filterClubTeams(teams, cfg.ClubName)
		if len(clubTeams) == 0 {
			fmt.Printf("    No teams matching \"%s\"\n", cfg.ClubName)
			continue
		}

		fmt.Printf("    Found %d matching teams: ", len(clubTeams))
		names := make([]string, len(clubTeams))
		for i, t := range clubTeams {
			names[i] = t.Label
		}
		fmt.Println(strings.Join(names, ", "))

		for _, team := range clubTeams {
			rounds, err := FetchFixture(comp.Value, section.Value, team.Value)
			if err != nil {
				fmt.Printf("    Error fetching fixture for %s: %v\n", team.Label, err)
				continue
			}

			data.Fixtures = append(data.Fixtures, models.TeamFixture{
				Team:    team,
				Section: section,
				Rounds:  rounds,
			})
		}
	}

	fmt.Printf("Scraping complete. Found %d team fixtures.\n\n", len(data.Fixtures))
	return data, nil
}

func filterClubTeams(teams []models.Team, clubName string) []models.Team {
	var matched []models.Team
	lower := strings.ToLower(clubName)
	for _, t := range teams {
		if strings.Contains(strings.ToLower(t.Label), lower) {
			matched = append(matched, t)
		}
	}
	return matched
}
