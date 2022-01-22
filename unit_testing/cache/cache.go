package cache

import (
	"fmt"
	"sync"
	"time"
)

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) (float64, error)
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	prices             map[string]CachedPrice
}

//	CachedPrice represents the value obtained from the service, denoting when it was retrieved
type CachedPrice struct {
	Value    float64
	DateFrom time.Time
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             map[string]CachedPrice{},
	}
}

var lock = sync.RWMutex{}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	lock.RLock()
	price, ok := c.prices[itemCode]
	lock.RUnlock()

	if ok {
		if time.Now().Sub(price.DateFrom) < c.maxAge {
			return price.Value, nil
		}
	}
	val, err := c.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, fmt.Errorf("getting price from service : %v", err.Error())
	}

	lock.Lock()
	c.prices[itemCode] = CachedPrice{
		Value:    val,
		DateFrom: time.Now(),
	}
	lock.Unlock()
	return val, nil
}

func (c *TransparentCache) GetConcurrentPriceFor(inputChannel chan string, itemAmount int) (chan float64, chan error) {
	resultChan := make(chan float64, itemAmount)
	errorChan := make(chan error, itemAmount)
	go func() {
		wg := &sync.WaitGroup{}
		//Changed infinite loop and tagged break for the range syntax
		// also changed the internal func sent to a goroutine to handle the itemCode by parameter
		for itemCode := range inputChannel {
			wg.Add(1)
			go func(code string) {
				defer (*wg).Done()
				price, err := c.GetPriceFor(code)
				if err != nil {
					errorChan <- err
				}
				resultChan <- price
			}(itemCode)
		}
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	return resultChan, errorChan
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	results := []float64{}

	var resultChan chan float64
	var errorChan chan error

	inputChanel := make(chan string, len(itemCodes))

	resultChan, errorChan = c.GetConcurrentPriceFor(inputChanel, len(itemCodes))

	for _, itemCode := range itemCodes {
		inputChanel <- itemCode
	}
	close(inputChanel)

	for {
		select {
		case result, isOpen := <-resultChan:
			if !isOpen {
				return results, nil
			}
			results = append(results, result)
		case err := <-errorChan:
			return results, err
		}
	}

}
