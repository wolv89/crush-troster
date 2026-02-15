package parser

import (
	"regexp"
	"strings"

	"github.com/wolv89/troster/internal/models"
)

var (
	selectOptionRe = regexp.MustCompile(`<option\s+value="([^"]*)"[^>]*>([^<]+)</option>`)
	fixtureRowRe   = regexp.MustCompile(`(?s)<tr>\s*<td[^>]*>([^<]*)</td>\s*<td[^>]*>([^<]*)</td>\s*<td[^>]*>(.*?)</td>\s*<td[^>]*>(.*?)</td>`)
	tagStripRe     = regexp.MustCompile(`<[^>]*>`)
	nbspRe         = regexp.MustCompile(`&nbsp;`)
	multiSpaceRe   = regexp.MustCompile(`\s+`)
)

func ParseCompetitions(html string) []models.Competition {
	return parseSelect(html, "daytime")
}

func ParseSections(html string) []models.Section {
	options := parseSelect(html, "section")
	sections := make([]models.Section, len(options))
	for i, o := range options {
		sections[i] = models.Section{Value: o.Value, Label: o.Label}
	}
	return sections
}

func ParseTeams(html string, sectionValue string) []models.Team {
	options := parseSelect(html, "team")
	teams := make([]models.Team, len(options))
	for i, o := range options {
		teams[i] = models.Team{Value: o.Value, Label: o.Label, Section: sectionValue}
	}
	return teams
}

func ParseFixture(html string, teamValue string) []models.FixtureRound {
	// Find the fixture table — it's the one with "Fixture for" heading
	fixtureIdx := strings.Index(html, "Fixture for")
	if fixtureIdx == -1 {
		return nil
	}

	// Extract from the fixture table onward
	fixtureHTML := html[fixtureIdx:]

	// Find all table rows with fixture data
	// The pattern: <tr><td>round</td><td>date</td><td>home</td><td>away</td>
	matches := fixtureRowRe.FindAllStringSubmatch(fixtureHTML, -1)

	var rounds []models.FixtureRound
	for _, m := range matches {
		round := cleanText(m[1])
		date := cleanText(m[2])
		homeRaw := m[3]
		awayRaw := m[4]

		if date == "" {
			continue
		}

		// Check for "No Play"
		if strings.Contains(homeRaw, "No Play") || strings.Contains(awayRaw, "No Play") {
			rounds = append(rounds, models.FixtureRound{
				Date:   date,
				NoPlay: true,
			})
			continue
		}

		home := cleanText(homeRaw)
		away := cleanText(awayRaw)

		if home == "" && away == "" {
			continue
		}

		rounds = append(rounds, models.FixtureRound{
			Round:    round,
			Date:     date,
			HomeTeam: home,
			AwayTeam: away,
		})
	}

	return rounds
}

func parseSelect(html string, selectName string) []models.Competition {
	// Find the select element by name
	selectPattern := regexp.MustCompile(`(?s)<select[^>]*name="` + selectName + `"[^>]*>(.*?)</select>`)
	selectMatch := selectPattern.FindStringSubmatch(html)
	if selectMatch == nil {
		return nil
	}

	selectHTML := selectMatch[1]
	matches := selectOptionRe.FindAllStringSubmatch(selectHTML, -1)

	var options []models.Competition
	for _, m := range matches {
		value := m[1]
		label := cleanText(m[2])
		if value == "" || label == "" {
			continue
		}
		options = append(options, models.Competition{Value: value, Label: label})
	}

	return options
}

func cleanText(s string) string {
	s = nbspRe.ReplaceAllString(s, " ")
	s = tagStripRe.ReplaceAllString(s, " ")
	s = multiSpaceRe.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}
