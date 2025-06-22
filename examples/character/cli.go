package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/karashiiro/bingode"
	"github.com/xivapi/godestone/v2"
)

// 保留原始的命令行工具功能
func runCLI() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: go run . <character_id>")
	}

	s := godestone.NewScraper(bingode.New(), godestone.EN)

	id, err := strconv.ParseUint(os.Args[1], 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	c, err := s.FetchCharacter(uint32(id))
	if err != nil {
		log.Fatalln(err)
	}

	cJSON, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(cJSON))
}