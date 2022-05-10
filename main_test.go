package main

import (
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func TestSync(t *testing.T) {
	main()
}
