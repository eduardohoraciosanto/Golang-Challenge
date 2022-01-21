package code_implementation

import (
	"errors"
	"sync"
)

type UserRepository interface {
	GetUser(id string) (User, error)
	InsertUser(userID string, friends []string) error
	FindAllSocialCircles(userIds []string) map[string][]string
}

type User struct {
	Id      string
	Friends []string
}

type userRepository struct {
	users map[string]User
	lock  *sync.RWMutex
}

func NewUsersRepository() UserRepository {
	return &userRepository{
		users: make(map[string]User),
		lock:  &sync.RWMutex{},
	}
}

// InsertUser allows for the insertion of a user into the known users.
//
// Returns error if user ID is already taken
func (u *userRepository) InsertUser(userID string, friends []string) error {

	//Lock for Reading
	u.lock.RLock()
	_, exists := u.users[userID]
	u.lock.RUnlock()

	if exists {
		return errors.New("UserID already exists")
	}
	//create new user
	user := User{
		Id:      userID,
		Friends: friends,
	}
	//Lock for Writing and Add user to the internal map
	u.lock.Lock()
	u.users[userID] = user
	u.lock.Unlock()
	return nil
}

// GetUser fetches a user from the known users.
//
// Returns error if not found.
// Returns the User and a nil error if found.
func (u *userRepository) GetUser(id string) (User, error) {
	return User{}, nil
}

//FindAllSocialCircles finds the social circles of the supplied userIds.
//
// userIds that are not in the known userbase won't exist in the returned map
func (u *userRepository) FindAllSocialCircles(userIds []string) map[string][]string {

	return nil
}

//getFriendsOf returns the friends for a particular user
func (u *userRepository) getFriendsOf(userID string) []string {

	return nil
}
