package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/mock"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create and start the mock server
	server := mock.NewMockServer(8083)
	log.Fatal(server.Start())
}
