package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/eduardohoraciosanto/Golang-Challenge/unit_testing/cache"
)

const (
	testDefaultMaxAge        = 3 * time.Second
	testDefaultResponseDelay = 200 * time.Millisecond
	timingRuleMetric         = 1500 * time.Millisecond
)

func TestGetPriceForOkNonCached(t *testing.T) {
	ps := &priceServiceMock{
		shouldFail:        false,
		priceReturn:       12.34,
		mockResponseDelay: testDefaultResponseDelay,
	}

	c := cache.NewTransparentCache(ps, testDefaultMaxAge)

	itemPrice, err := c.GetPriceFor("someItemCode")
	if err != nil {
		t.Fatalf("GetPriceFor not expected to fail, err: %+v", err)
	}
	if itemPrice != ps.priceReturn {
		t.Fatalf("Unexpected price returned. Expected: %v, Got: %v", ps.priceReturn, itemPrice)
	}

}

func TestGetPriceForOkCached(t *testing.T) {
	ps := &priceServiceMock{
		shouldFail:        false,
		priceReturn:       12.34,
		mockResponseDelay: testDefaultResponseDelay,
	}

	c := cache.NewTransparentCache(ps, testDefaultMaxAge)

	//first call will cache the item price
	_, err := c.GetPriceFor("someItemCode")
	if err != nil {
		t.Fatalf("GetPriceFor not expected to fail, err: %+v", err)
	}

	//second call will give us the cached priced
	itemPrice, err := c.GetPriceFor("someItemCode")
	if err != nil {
		t.Fatalf("GetPriceFor not expected to fail, err: %+v", err)
	}

	if itemPrice != ps.priceReturn {
		t.Fatalf("Unexpected price returned. Expected: %v, Got: %v", ps.priceReturn, itemPrice)
	}

}

func TestGetPriceForError(t *testing.T) {
	ps := &priceServiceMock{
		shouldFail:        true,
		priceReturn:       12.34,
		mockResponseDelay: testDefaultResponseDelay,
	}

	c := cache.NewTransparentCache(ps, testDefaultMaxAge)

	_, err := c.GetPriceFor("someItemCode")
	if err == nil {
		t.Fatalf("GetPriceFor expected to fail")
	}

}

func TestGetPricesForOK(t *testing.T) {
	ps := &priceServiceMock{
		shouldFail:        false,
		priceReturn:       12.34,
		mockResponseDelay: 500 * time.Millisecond,
	}
	items := []string{"SomeItemCode", "AnotherItemCode", "CoolItemCode"}
	c := cache.NewTransparentCache(ps, testDefaultMaxAge)

	_, err := c.GetPricesFor(items...)
	if err != nil {
		t.Fatalf("GetPricesFor was not expected to fail. Error: %+v", err)
	}
}

func TestGetPricesForError(t *testing.T) {
	ps := &priceServiceMock{
		shouldFail:        true,
		priceReturn:       12.34,
		mockResponseDelay: 500 * time.Millisecond,
	}
	items := []string{"SomeItemCode", "AnotherItemCode", "CoolItemCode"}
	c := cache.NewTransparentCache(ps, testDefaultMaxAge)

	_, err := c.GetPricesFor(items...)
	if err == nil {
		t.Fatalf("GetPricesFor should have failed")
	}
}

func TestGetPricesTimingRule(t *testing.T) {
	ps := &priceServiceMock{
		shouldFail:        false,
		priceReturn:       12.34,
		mockResponseDelay: 1 * time.Second,
	}
	items := []string{"Item1", "Item2", "Item3", "Item4", "Item5"}
	c := cache.NewTransparentCache(ps, testDefaultMaxAge)

	t0 := time.Now()
	_, err := c.GetPricesFor(items...)
	if err != nil {
		t.Fatalf("GetPricesFor was not expected to fail. Error: %+v", err)
	}
	tdelta := time.Since(t0)

	if tdelta >= timingRuleMetric {
		t.Fatalf("GetPricesFor 5 items should take less than %v, it took %v", timingRuleMetric, tdelta)
	}
}

//mocks
type priceServiceMock struct {
	shouldFail        bool
	priceReturn       float64
	mockResponseDelay time.Duration
}

func (ps *priceServiceMock) GetPriceFor(itemCode string) (float64, error) {
	if ps.shouldFail {
		return 0, errors.New("Mock was asked to fail")
	}
	time.Sleep(ps.mockResponseDelay)
	return ps.priceReturn, nil
}
