package main

import (
	"github.com/c-mierez/rss-aggregator/internal/env"
)

func main() {
	// Check Environment Variable definitions
	env.LoadAndCheckENV(true)

}
