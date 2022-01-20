# Golang-Challenge

We ask that you complete the following challenge to evaluate your development skills. Please use the programming language and framework discussed during your interview to accomplish the following task.

PLEASE DO NOT FORK THE REPOSITORY. WE NEED A PUBLIC REPOSITORY FOR THE REVIEW.

## Show your work

Create a Public repository ( please dont make a pull request, clone the private repository and create a new plublic one on your profile)
Commit each step of your process so we can follow your thought process.

# Code Implementation

Given a list of userIds, you have to find the social circle of each one of them. A social circle of a user consists of their friends and the friends of their friends, and so on. Let’s see an example:

    UserA: [B]
    UserB: [A]
    UserC: [D, F]
    UserD: [C]
    UserE: []
    UserF: [C, G]
    UserG: [F]
    Graph representing the friendships between users:
            A
           /
          B      C
                / |
          E    D  F
                  |    
                  G


Social circle of A --> [B]

Social circle of D --> [C, F, G]

## Challenge
Develop the body of the following function

    func FindAllSocialCircles(userIds []string) map[string][]string {}

**Input:** an array of userIds

**Output:** a map which each key-value pair represent userId -> [socialCircle]

In the above example, the execution would be the following:

    result := FindAllSocialCircles([A,D])

And result is a map:

    {
        “A”: [B],
        “D”: [C, F, G]
    }

Here are 2 types which you have to use as a baseline:

    type UserRepository interface {
        GetUser(id string) (User, error)
    }

    type User struct {
        Id  	string
        Friends []string
    }

What we would like to see in your solution:

* Use the concurrency tools provided by the language.
* Think of a solution which performs well for large datasets.
* Use exit conditionals.
* Add one unit test using the provided example data to validate the code.
* Consider that the exercise is about friends, and human ralationships are bidirectional

### Note:
You can use UserRepository as a global variable for simplicity. This interface represents fetching a user from an external service.

## Time Spent
We suggest not spending more than 4-5 hours total. Please make commits as often as possible so we can see the time you spent and please do not make one commit. We will evaluate the code and time spent.

What we want to see is how well you handle yourself given the time you spend on the problem, how you think, and how you prioritize when time is sufficient to solve everything.

Please email your solution as soon as you have completed the challenge or the time is up
