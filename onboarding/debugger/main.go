package main

import (
	"onboarding/api/api"
	"onboarding/guessers/guessers"
	"onboarding/numbers/numbers"
)

func main() {
	go api.RealApi()
	go numbers.RealNumbers()
	go guessers.RealGuessers()
	select {}
}
