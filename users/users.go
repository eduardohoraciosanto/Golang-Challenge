package users

import (
	"errors"
)

type UserRepository interface {
	GetUser(id string) (User, error)
}

type User struct {
	Id      string
	Friends []string
}

type userRepository struct {
}

var (
	ErrorUserAlreadyExists = errors.New("UserID already exists")
	ErrorUserNotFound      = errors.New("UserID not found")
	userRepo               = userRepository{}
)

// GetUser fetches a user from the known users.
//
// Returns error if not found.
// Returns the User and a nil error if found.
func (u *userRepository) GetUser(id string) (User, error) {
	switch id {
	case "A":
		return User{
			Id:      "A",
			Friends: []string{"B"},
		}, nil
	case "B":
		return User{
			Id:      "B",
			Friends: []string{"A"},
		}, nil
	case "C":
		return User{
			Id:      "C",
			Friends: []string{"D", "F"},
		}, nil
	case "D":
		return User{
			Id:      "D",
			Friends: []string{"C"},
		}, nil
	case "E":
		return User{
			Id:      "E",
			Friends: []string{},
		}, nil
	case "F":
		return User{
			Id:      "F",
			Friends: []string{"C", "G"},
		}, nil
	case "G":
		return User{
			Id:      "G",
			Friends: []string{"F"},
		}, nil
	}

	return User{}, ErrorUserNotFound
}

//FindAllSocialCircles finds the social circles of the supplied userIds.
//
// userIds that are not in the known userbase won't exist in the returned map
func FindAllSocialCircles(userIds []string) map[string][]string {

	//Code Red, then I'll make it Green
	for _, id := range userIds {
		//I'll use a map with an empty struct to save a bit on memory and avoid repetitions
		sc := make(map[string]struct{})

		user, err := userRepo.GetUser(id)
		if err != nil {
			//user that does not exist
			continue
		}
		for _, friendID := range user.Friends {
			sc[friendID] = struct{}{}
		}

	}

	return nil
}
