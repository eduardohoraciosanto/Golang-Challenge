# Unit Testing

The purpose of this challenge is creating Unit tests for the given working piece of code. Please feel free to use the techniques you consider more suitable and any code refactor needed to make sure the best testing is achieved.

## Challenge
This code implements a simple cache for a service that consumes a lot of time for each request that’s made against it. The solution provided uses an in memory map to store already consulted values, and concurrently searches for all the items asked to find.

You are asked to test the cache.go file.

Things to bear in mind:
* A minimum of 90% coverage is required
* **Concurrency:** Considering that PriceServices takes a full second for each call (this should be mocked), test that the function takes no more than 1.5 seconds for 5 different items.


### Note:
If you think that a part of the solution provided needs to be refactored, feel free to do it (and if you do, please leave a comment, so we don’t forget to check). Remember this challenge is focused on unit testing, so do not dedicate a lot of time refactoring.