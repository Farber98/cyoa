package main

import (
	"cyoa"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	jsonFlag := flag.String("json", "../gopher.json", ".json filename to read from")
	flag.Parse()

	jsonFile, err := os.Open(*jsonFlag)
	if err != nil {
		log.Panic("jsonFile parsing err   #%v ", err)
	}

	story, err := cyoa.StoryToJson(jsonFile)
	if err != nil {
		log.Panic("Error parsing json story.")
	}
	fmt.Println(story)
}
