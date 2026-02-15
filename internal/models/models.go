package models

type Config struct {
	ClubName     string `json:"club_name"`
	TotalCourts  int    `json:"total_courts"`
	CourtsPerTeam int   `json:"courts_per_team"`
}

type Competition struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Section struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Team struct {
	Value   string `json:"value"`
	Label   string `json:"label"`
	Section string `json:"section"`
}

type FixtureRound struct {
	Round    string `json:"round"`
	Date     string `json:"date"`
	HomeTeam string `json:"home_team"`
	AwayTeam string `json:"away_team"`
	IsHome   bool   `json:"is_home"`
	NoPlay   bool   `json:"no_play"`
}

type TeamFixture struct {
	Team     Team           `json:"team"`
	Section  Section        `json:"section"`
	Rounds   []FixtureRound `json:"rounds"`
}

type ScrapedData struct {
	Competition Competition   `json:"competition"`
	Config      Config        `json:"config"`
	Fixtures    []TeamFixture `json:"fixtures"`
	ScrapedAt   string        `json:"scraped_at"`
}
