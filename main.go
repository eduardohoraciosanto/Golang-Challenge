package main

import (
	"fmt"
	"time"

	"github.com/eduardohoraciosanto/Golang-Challenge/users"
)

// Here I'll make use of the UserRepository interface, adding users and then finding the social circle for them

func main() {
	now := time.Now()
	socialCircles := users.FindAllSocialCircles([]string{"A", "D", "G"})
	took := time.Since(now)
	fmt.Printf("Social Cirles: %+v. We took: %v", socialCircles, took)
}
