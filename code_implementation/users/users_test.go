package users_test

import (
	"testing"

	"github.com/eduardohoraciosanto/Golang-Challenge/code_implementation/users"
	"gotest.tools/assert"
)

func TestFindAllSocialCirclesOK(t *testing.T) {
	testUsers := []string{"A", "D"}
	expected := map[string][]string{
		"A": []string{"B"},
		"D": []string{"C", "F", "G"},
	}

	result := users.FindAllSocialCircles(testUsers)

	assert.DeepEqual(t, expected, result)

}
