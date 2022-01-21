package main

import (
	"fmt"
	"time"

	"github.com/eduardohoraciosanto/Golang-Challenge/users"
)

func main() {
	now := time.Now()
	socialCircles := users.FindAllSocialCircles([]string{"A", "D", "G"})
	took := time.Since(now)
	fmt.Printf("Social Cirles: %+v. We took: %v", socialCircles, took)
}
