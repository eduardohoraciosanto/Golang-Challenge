package users

import (
	"errors"
	"sort"
	"sync"
)

type UserRepository interface {
	GetUser(id string) (User, error)
}

type User struct {
	Id      string
	Friends []string
}

//userRepository is just a struct to allow the implementation of the UserRepository
type userRepository struct {
}

//ConcurrentMap allows for using a Map safely
type ConcurrentMap interface {
	FriendExists(user, friend string) bool
	AddFriend(user, friend string)
	GetSocialCircles() map[string][]string
}

type concurrentMap struct {
	data map[string]map[string]struct{}
	lock *sync.RWMutex
}

//FriendExists is used to safely check for an existing asociation userID<-> friend
func (c *concurrentMap) FriendExists(user, friend string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	friends, exists := c.data[user]
	if !exists {
		return false
	}
	_, friendExists := friends[friend]

	return friendExists
}

//AddFriend allows for adding a friend safely on the concurrent map
func (c *concurrentMap) AddFriend(user, friend string) {
	//we won't add the user as self-friend
	if user == friend {
		return
	}
	//friends already existing won't be added
	if c.FriendExists(user, friend) {
		return
	}

	//safe access to underlying map
	c.lock.Lock()
	defer c.lock.Unlock()

	//get friends, add new one, save it.
	friends := c.data[user]
	if friends == nil {
		friends = make(map[string]struct{})
	}
	friends[friend] = struct{}{}
	c.data[user] = friends
}

//GetSocialCircles returns the social circles in the necessary format
func (c *concurrentMap) GetSocialCircles() map[string][]string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	socialCircles := make(map[string][]string)

	for user, userFriends := range c.data {
		friends := []string{}
		for friend := range userFriends {
			friends = append(friends, friend)
		}
		sort.SliceStable(friends, func(i, j int) bool {
			return friends[i] < friends[j]
		})
		socialCircles[user] = friends
	}
	return socialCircles
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
	wg := &sync.WaitGroup{}
	socialCircles := &concurrentMap{
		data: make(map[string]map[string]struct{}),
		lock: &sync.RWMutex{},
	}
	// using two goroutines in parallel for each userID.
	// synchronized them concurrently using channels. one for each userID
	for _, id := range userIds {
		friendsChan := make(chan []string)
		go findAllFriends(id, friendsChan)
		wg.Add(1)
		go gatherFriends(id, friendsChan, socialCircles, wg)
	}
	wg.Wait()
	return socialCircles.GetSocialCircles()
}

//findAllFriends traverses the friends of the friends and sends information on the data channel. It also closes the channel once the recursive func has ended
//
// To be used in parallel with gatherFriends
func findAllFriends(userID string, dataChan chan []string) {
	recursiveFriendsOf(userID, nil, dataChan)
	close(dataChan)
}

//recursiveFriendsOf keeps track of what friends were already sent to the channel and keeps sending data unless all friends are sent
func recursiveFriendsOf(userID string, friendsOfFound map[string]struct{}, dataChan chan []string) {
	if friendsOfFound == nil {
		friendsOfFound = make(map[string]struct{})
	}

	if _, found := friendsOfFound[userID]; found {
		return
	}

	user, err := userRepo.GetUser(userID)
	if err != nil {
		return
	}
	dataChan <- user.Friends
	friendsOfFound[userID] = struct{}{}
	for _, friendID := range user.Friends {
		recursiveFriendsOf(friendID, friendsOfFound, dataChan)
	}
}

//gatherFriends receives friends over the channel and adds them safely to the underlying safeMap.
//
// To be used in parallel with findAllFriends
func gatherFriends(userID string, datachan chan []string, friendsMap ConcurrentMap, wg *sync.WaitGroup) {
	for incomingFriends := range datachan {
		for _, friend := range incomingFriends {
			friendsMap.AddFriend(userID, friend)
		}
	}
	wg.Done()
}
