package main

import (
	"log"
	config2 "studyAIModel/config"
)

func main() {
	config, err := config2.Load()
	if err != nil {
		log.Fatal(err)
	}
}
