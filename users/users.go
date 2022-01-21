package users

import (
	"errors"
	"sync"
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

type ConcurrentMap interface {
	Exists(id string) bool
	Add(id string)
	Remove(id string)
}
type concurrentMap struct {
	data map[string]struct{}
	lock *sync.RWMutex
}

func (c *concurrentMap) Exists(id string) bool {
	c.lock.RLock()
	_, exists := c.data[id]
	c.lock.RUnlock()
	return exists
}

func (c *concurrentMap) Add(id string) {
	c.lock.Lock()
	c.data[id] = struct{}{}
	c.lock.Unlock()
}

func (c *concurrentMap) Remove(id string) {
	c.lock.Lock()
	delete(c.data, id)
	c.lock.Unlock()
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
	circles := make(map[string][]string)
	//Code Red, then I'll make it Green
	for _, id := range userIds {
		socialCircle := concurrentMap{
			data: make(map[string]struct{}),
			lock: &sync.RWMutex{},
		}
		findFriends(id, &socialCircle)
		socialCircle.Remove(id)

		friendsSlice := []string{}
		for friendID := range socialCircle.data {
			friendsSlice = append(friendsSlice, friendID)
		}
		circles[id] = friendsSlice
	}

	return circles
}

func findFriends(userID string, dataMap ConcurrentMap) {
	user, err := userRepo.GetUser(userID)
	if err != nil {
		//user that does not exist
		return
	}
	for _, friendID := range user.Friends {
		if dataMap.Exists(friendID) {
			continue
		}
		dataMap.Add(friendID)
		findFriends(friendID, dataMap)
	}

}
