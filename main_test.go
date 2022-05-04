package main

import (
	"github.com/joho/godotenv"
	"testing"
)

func init() {
	godotenv.Load()
}

func TestSync(t *testing.T) {
	main()
}
