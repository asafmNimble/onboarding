package main

import (
	"fmt"
	"onboarding/api/api"
	"onboarding/guessers/guessers"
	"onboarding/numbers/numbers"
	"time"
)

func main() {
	fmt.Println("\n\n", time.Unix(1642509050, 0).UTC(), "\n\n")
	go api.RealApi()
	go numbers.RealNumbers()
	go guessers.RealGuessers()
	select {}
}
