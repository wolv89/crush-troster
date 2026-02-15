package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wolv89/troster/internal/models"
)

type Flags struct {
	ClubName     string
	TotalCourts  int
	CourtsPerTeam int
	ForceScrape  bool
}

func ParseFlags() Flags {
	f := Flags{}
	flag.StringVar(&f.ClubName, "club", "", "Tennis club name (e.g. \"Hurlingham Park\")")
	flag.IntVar(&f.TotalCourts, "courts", 0, "Total number of courts at the club")
	flag.IntVar(&f.CourtsPerTeam, "per-team", 0, "Number of courts needed per home team")
	flag.BoolVar(&f.ForceScrape, "force-scrape", false, "Force re-scrape even if cached data exists")
	flag.Parse()
	return f
}

func GetConfig(f Flags) models.Config {
	reader := bufio.NewReader(os.Stdin)
	cfg := models.Config{
		ClubName:     f.ClubName,
		TotalCourts:  f.TotalCourts,
		CourtsPerTeam: f.CourtsPerTeam,
	}

	if cfg.ClubName == "" {
		cfg.ClubName = prompt(reader, "Enter your tennis club name")
	}
	if cfg.TotalCourts == 0 {
		cfg.TotalCourts = promptInt(reader, "How many courts does your club have")
	}
	if cfg.CourtsPerTeam == 0 {
		cfg.CourtsPerTeam = promptInt(reader, "How many courts are needed per home team")
	}

	fmt.Printf("\nClub: %s | Courts: %d | Per team: %d\n\n", cfg.ClubName, cfg.TotalCourts, cfg.CourtsPerTeam)
	return cfg
}

func ChooseCompetition(comps []models.Competition) models.Competition {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Available competitions:")
	for i, c := range comps {
		fmt.Printf("  %d. %s\n", i+1, c.Label)
	}
	fmt.Println()

	for {
		choice := promptInt(reader, "Select a competition (number)")
		if choice >= 1 && choice <= len(comps) {
			selected := comps[choice-1]
			fmt.Printf("\nSelected: %s\n\n", selected.Label)
			return selected
		}
		fmt.Printf("Please enter a number between 1 and %d\n", len(comps))
	}
}

func prompt(reader *bufio.Reader, label string) string {
	for {
		fmt.Printf("%s: ", label)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
	}
}

func promptInt(reader *bufio.Reader, label string) int {
	for {
		s := prompt(reader, label)
		n, err := strconv.Atoi(s)
		if err == nil && n > 0 {
			return n
		}
		fmt.Println("Please enter a valid positive number")
	}
}
