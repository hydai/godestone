# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the **godestone** library - a Go scraper for The Lodestone (Final Fantasy XIV's official website). The library provides functionality to scrape character profiles, achievements, minions, mounts, free companies, linkshells, and other game-related data.

### Architecture

- **Core Scraper**: The main `Scraper` struct in `scraper.go` handles all web scraping operations using the colly framework
- **Data Models**: Character, FreeCompany, Linkshell, CWLS, PVPTeam models defined in separate `.go` files
- **CSS Selectors**: Stored in `internal/selectors/` and packed using go-bindata via `generate.sh`
- **Data Provider**: Uses the `bingode` package for game data lookup
- **Examples**: Located in `examples/` directory showing basic usage patterns

The current working directory is in `examples/character/` which contains a simple example (`main.go`) that fetches and displays character data in JSON format.

## Development Commands

### Building and Running
```bash
# Build the example
go build main.go

# Run the character example (requires character ID argument)
go run main.go <character_id>

# Build the entire project from root
go build ./...
```

### Testing
```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run specific test
go test -run TestFetchCharacter
```

### CSS Selector Generation
```bash
# Repack CSS selectors (from project root)
./generate.sh
```

### Dependencies
- Uses Go modules (`go.mod`)
- Requires `go-bindata` for CSS selector packing
- Uses `github.com/karashiiro/bingode` for game data
- Built on `github.com/gocolly/colly/v2` for web scraping

## Key Components

- **Scraper Creation**: Always initialize with `NewScraper(dataProvider, lang)` 
- **Language Support**: EN, JA, FR, DE supported (not "zh")
- **Concurrent Design**: Many functions use channels for streaming results
- **Error Handling**: Most scraping methods return errors for network/parsing failures
- **Rate Limiting**: Built into colly framework used by the scraper

The library is designed for scraping The Lodestone website and requires a data provider (like bingode) to function properly.