package main

import (
	"fmt"

	"github.com/eduardohoraciosanto/Golang-Challenge/users"
)

// Here I'll make use of the UserRepository interface, adding users and then finding the social circle for them

func main() {

	socialCircles := users.FindAllSocialCircles([]string{"A", "F"})

	fmt.Printf("Social Cirles: %+v", socialCircles)
}
