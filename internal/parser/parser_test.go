package parser

import (
	"testing"
)

const competitionsHTML = `<select id="daytime" name="daytime" Onchange="select_submit(select,0);">
<option value="">&nbsp;</option>
<option value="AA">Saturday AM - Autumn 2026</option>
<option value="AP">Saturday PM - Summer 2025/26</option>
<option value="UA">Sunday AM - Autumn 2026</option>
</select>`

const sectionsHTML = `<select id="section" name="section" Onchange="select_submit(select,1);">
<option value="">&nbsp;</option>
<option value="AA001">Rubbers 1</option>
<option value="AA002">Sets 1</option>
</select>`

const teamsHTML = `<select id="team" name="team" Onchange="select_submit(select,2);">
<option value="">&nbsp;</option>
<option value="AA001">BLTC</option>
<option value="AA002">Clarinda</option>
<option value="AA003">Hurlingham Park Hammers</option>
<option value="AA004">Hurlingham Park Battlers</option>
</select>`

const fixtureHTML = `<td colspan="4" class="mg">Fixture for BLTC </td></tr>
<tr><th>Rd</th><th>Date</th><th style="text-align:left">Home</th><th style="text-align:left">Away</th></tr>
<tr><td style="text-align:center">1</td><td style="text-align:center">31 Jan 26</td><td>Clarinda&nbsp;<span style="font-style:italic">(10:00)</span></td><td>BLTC</td><td>&nbsp;<br/>&nbsp;</td></tr>
<tr><td style="text-align:center">2</td><td style="text-align:center">7 Feb 26</td><td>BLTC&nbsp;<span style="font-style:italic">(8:00)</span></td><td>Port Melbourne</td><td>&nbsp;<br/>&nbsp;</td></tr>
<tr><td style="text-align:center">&nbsp;</td><td style="text-align:center">7 Mar 26</td><td colspan="2"><b>No Play</b></td><td>&nbsp;<br/>&nbsp;</td></tr>
<tr><td style="text-align:center">6</td><td style="text-align:center">14 Mar 26</td><td>BLTC&nbsp;<span style="font-style:italic">(8:00)</span></td><td>Parkdale</td><td>&nbsp;<br/>&nbsp;</td></tr>`

func TestParseCompetitions(t *testing.T) {
	comps := ParseCompetitions(competitionsHTML)
	if len(comps) != 3 {
		t.Fatalf("expected 3 competitions, got %d", len(comps))
	}
	if comps[0].Value != "AA" {
		t.Errorf("expected first value AA, got %s", comps[0].Value)
	}
	if comps[0].Label != "Saturday AM - Autumn 2026" {
		t.Errorf("expected first label 'Saturday AM - Autumn 2026', got %s", comps[0].Label)
	}
}

func TestParseSections(t *testing.T) {
	sections := ParseSections(sectionsHTML)
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}
	if sections[0].Value != "AA001" {
		t.Errorf("expected first value AA001, got %s", sections[0].Value)
	}
}

func TestParseTeams(t *testing.T) {
	teams := ParseTeams(teamsHTML, "AA001")
	if len(teams) != 4 {
		t.Fatalf("expected 4 teams, got %d", len(teams))
	}
	if teams[2].Label != "Hurlingham Park Hammers" {
		t.Errorf("expected 'Hurlingham Park Hammers', got %s", teams[2].Label)
	}
	if teams[2].Section != "AA001" {
		t.Errorf("expected section AA001, got %s", teams[2].Section)
	}
}

func TestParseFixture(t *testing.T) {
	rounds := ParseFixture(fixtureHTML, "AA001")
	if len(rounds) != 4 {
		t.Fatalf("expected 4 rounds, got %d", len(rounds))
	}

	if rounds[0].Round != "1" || rounds[0].Date != "31 Jan 26" {
		t.Errorf("round 1: got round=%q date=%q", rounds[0].Round, rounds[0].Date)
	}
	if rounds[0].HomeTeam != "Clarinda (10:00)" {
		t.Errorf("round 1 home: got %q", rounds[0].HomeTeam)
	}

	if !rounds[2].NoPlay {
		t.Error("expected round 3 to be NoPlay")
	}

	if rounds[3].Round != "6" || rounds[3].HomeTeam != "BLTC (8:00)" {
		t.Errorf("round 4: got round=%q home=%q", rounds[3].Round, rounds[3].HomeTeam)
	}
}
