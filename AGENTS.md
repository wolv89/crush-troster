# Agents Guide — Troster

Tennis fixture roster app. Scrapes fixtures from [TROLS](https://trols.org.au/brta/fixture.php), serves a local web UI to view home/away schedules.

## Commands

```bash
# Build
GOMODCACHE=$HOME/.gomodcache go build -o troster .

# Run (interactive prompts)
./troster

# Run (with flags, skip prompts)
./troster -club "Hurlingham Park" -courts 6 -per-team 2

# Force re-scrape (ignore cache)
./troster -force-scrape

# Tests
GOMODCACHE=$HOME/.gomodcache go test ./...

# Vet
GOMODCACHE=$HOME/.gomodcache go vet ./...
```

**Note**: The default `GOMODCACHE` (`/go/pkg/mod/cache`) is read-only in this distrobox. Always set `GOMODCACHE=$HOME/.gomodcache` for build/test commands.

## Project Structure

```
main.go                      # Entry point: CLI → scrape → serve
internal/
  cli/cli.go                 # Flag parsing + interactive prompts
  scraper/scraper.go         # HTTP POST to TROLS, orchestrates full scrape
  parser/parser.go           # Regex-based HTML parsing (select options, fixture tables)
  parser/parser_test.go      # Unit tests for all parse functions
  cache/cache.go             # JSON file cache with 24h TTL
  models/models.go           # Shared types (Config, Competition, Section, Team, etc.)
  server/server.go           # HTTP server: /api/fixtures JSON + static file serving
web/public/
  index.html                 # React SPA (CDN-loaded, Babel-compiled in-browser)
data/                        # Cache directory (gitignored, created at runtime)
```

## Architecture

### Flow
1. CLI collects config (club name, courts, courts per team) via flags or interactive prompts
2. Scraper fetches competition list from TROLS
3. User selects competition in terminal
4. Scraper walks all sections → finds club teams by name match → fetches each team's fixture
5. Data cached as JSON in `data/fixtures_{comp}.json`
6. HTTP server starts on `:8080`, serves React frontend + `/api/fixtures` endpoint

### TROLS Scraping Protocol
The TROLS site uses form POST navigation with three levels:
1. **Competition**: POST `daytime=XX` → returns sections dropdown
2. **Section**: POST `daytime=XX&section=XXXXX&which=1` → returns teams dropdown
3. **Team fixture**: POST `daytime=XX&section=XXXXX&team=XXXXX&which=2` → returns fixture table

- Form action: `https://trols.org.au/brta/fixture.php`
- Content-Type: `application/x-www-form-urlencoded`
- Hidden fields: `which` (navigation level), `style` (always empty)
- 500ms delay between requests to be polite to the server

### Fixture Table Structure
Each fixture row has: Round number, Date, Home team (with time), Away team. "No Play" rows have no round number and `<b>No Play</b>` spanning the home/away columns. Team names may include HTML spans for color (e.g., `Ormond <span style="color:red">Red</span>`).

## Conventions

- **Go standard library only** — no external dependencies
- **Frontend**: Single HTML file with React 18 via CDN + Babel standalone (no build step, no Node.js)
- **Caching**: JSON files in `data/`, 24h TTL, `--force-scrape` to bypass
- **HTML parsing**: Regex-based (not a DOM parser) since TROLS HTML is simple table-based
- **Team matching**: Case-insensitive `strings.Contains` on team name vs club name
- **One competition at a time**: relaunch to switch competitions

## Gotchas

- TROLS is a legacy PHP site with HTML table layouts. The form uses `onchange` JS to auto-submit — we replicate this with direct POST requests
- The `which` hidden field controls navigation depth (0=competition, 1=section, 2=team)
- Some teams have colored suffixes (Red, Blue, Yellow) rendered as inline `<span>` tags — the parser strips these to plain text
- "No Play" rows lack a round number; the round `<td>` contains only `&nbsp;`
- The fixture table for a team shows ALL rounds including opponent fixtures — we parse them all and the frontend determines home/away by matching against the club name
- The scraper filters teams at the dropdown level (before fetching fixtures) using the club name, so irrelevant teams are never fetched
- React is compiled in-browser via Babel standalone — no build tooling required, but browser must load CDN scripts

## Future Work (from PLAN.md)

Court assignment/rostering: assign teams to specific courts, shuffle, track balance of court usage across home and away matches. The frontend is structured with this in mind (tab-based views, config bar showing court info).
