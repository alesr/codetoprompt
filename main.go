package main

import (
	"log"

	"github.com/alesr/codetoprompt/internal/codetoprompt"
)

func main() {
	if err := codetoprompt.Run(); err != nil {
		log.Fatalln(err)
	}
}
