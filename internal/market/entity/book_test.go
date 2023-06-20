package entity

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuyAsset(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor1 := NewInvestor("1")
	investor2 := NewInvestor("2")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 10)
	investor1.AddAssetPosition(investorAssetPosition)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	sellOrder := NewOrder("1", investor1, asset1, 5, 5, "SELL")
	orderChan <- sellOrder

	buyOrder := NewOrder("2", investor2, asset1, 5, 5, "BUY")
	orderChan <- buyOrder
	wg.Wait()

	assert := assert.New(t)
	assert.Equal("CLOSED", sellOrder.Status, "Order 1 should be closed")
	assert.Equal(0, sellOrder.PendingShares, "Order 1 should have 0 PendingShares")
	assert.Equal("CLOSED", buyOrder.Status, "Order 2 should be closed")
	assert.Equal(0, buyOrder.PendingShares, "Order 2 should have 0 PendingShares")

	assert.Equal(5, investorAssetPosition.Shares, "Investor 1 should have 5 shares of asset 1")
	assert.Equal(5, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 5 shares of asset 1")
}

func TestBuyAssetWithDifferentAssents(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)
	asset2 := NewAsset("asset2", "Asset 2", 100)

	investor := NewInvestor("1")
	investor2 := NewInvestor("2")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 10)
	investor.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewInvestorAssetPosition("asset2", 10)
	investor2.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	order := NewOrder("1", investor, asset1, 5, 5, "SELL")
	orderChan <- order

	order2 := NewOrder("2", investor2, asset2, 5, 5, "BUY")
	orderChan <- order2

	assert := assert.New(t)
	assert.Equal("OPEN", order.Status, "Order 1 should be closed")
	assert.Equal(5, order.PendingShares, "Order 1 should have 5 PendingShares")
	assert.Equal("OPEN", order2.Status, "Order 2 should be closed")
	assert.Equal(5, order2.PendingShares, "Order 2 should have 5 PendingShares")
}

func TestBuyPartialAsset(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor1 := NewInvestor("1")
	investor2 := NewInvestor("2")
	investor3 := NewInvestor("3")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 3)
	investor1.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewInvestorAssetPosition("asset1", 5)
	investor3.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	// investidor 2 quer comprar 5 shares
	buyOrder := NewOrder("1", investor2, asset1, 5, 5.0, "BUY")
	orderChan <- buyOrder

	// investidor 1 quer vender 3 shares
	sellOrder := NewOrder("2", investor1, asset1, 3, 5.0, "SELL")
	orderChan <- sellOrder

	assert := assert.New(t)
	go func() {
		for range orderChanOut {
		}
	}()

	wg.Wait()

	assert.Equal("CLOSED", sellOrder.Status, "Order 1 should be closed")
	assert.Equal(0, sellOrder.PendingShares, "Order 1 should have 0 PendingShares")

	assert.Equal("OPEN", buyOrder.Status, "Order 2 should be OPEN")
	assert.Equal(2, buyOrder.PendingShares, "Order 2 should have 2 PendingShares")

	assert.Equal(0, investorAssetPosition.Shares, "Investor 1 should have 0 shares of asset 1")
	assert.Equal(3, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 3 shares of asset 1")

	wg.Add(1)
	order3 := NewOrder("3", investor3, asset1, 2, 5.0, "SELL")
	orderChan <- order3
	wg.Wait()

	assert.Equal("CLOSED", order3.Status, "Order 3 should be closed")
	assert.Equal(0, order3.PendingShares, "Order 3 should have 0 PendingShares")

	assert.Equal("CLOSED", buyOrder.Status, "Order 2 should be CLOSED")
	assert.Equal(0, buyOrder.PendingShares, "Order 2 should have 0 PendingShares")

	assert.Equal(2, len(book.Transactions), "Should have 2 transactions")
	assert.Equal(15.0, float64(book.Transactions[0].Total), "Transaction should have price 15")
	assert.Equal(10.0, float64(book.Transactions[1].Total), "Transaction should have price 10")
}

func TestBuyWithDifferentPrice(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor1 := NewInvestor("1")
	investor2 := NewInvestor("2")
	investor3 := NewInvestor("3")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 3)
	investor1.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewInvestorAssetPosition("asset1", 5)
	investor3.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	// investidor 2 quer comprar 5 shares
	buyOrder := NewOrder("2", investor2, asset1, 5, 5.0, "BUY")
	orderChan <- buyOrder

	// investidor 1 quer vender 3 shares
	sellOrder := NewOrder("1", investor1, asset1, 3, 4.0, "SELL")
	orderChan <- sellOrder

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	assert := assert.New(t)
	assert.Equal("CLOSED", sellOrder.Status, "Order 1 should be closed")
	assert.Equal(0, sellOrder.PendingShares, "Order 1 should have 0 PendingShares")

	assert.Equal("OPEN", buyOrder.Status, "Order 2 should be OPEN")
	assert.Equal(2, buyOrder.PendingShares, "Order 2 should have 2 PendingShares")

	assert.Equal(0, investorAssetPosition.Shares, "Investor 1 should have 0 shares of asset 1")
	assert.Equal(3, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 3 shares of asset 1")

	wg.Add(1)
	sellOrder2 := NewOrder("3", investor3, asset1, 3, 4.5, "SELL")
	orderChan <- sellOrder2

	wg.Wait()

	assert.Equal("OPEN", sellOrder2.Status, "Order 3 should be open")
	assert.Equal(1, sellOrder2.PendingShares, "Order 3 should have 1 PendingShares")

	assert.Equal("CLOSED", buyOrder.Status, "Order 2 should be CLOSED")
	assert.Equal(0, buyOrder.PendingShares, "Order 2 should have 0 PendingShares")
}

func TestNoMatch(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor := NewInvestor("1")
	investor2 := NewInvestor("2")

	investorAssetPosition := NewInvestorAssetPosition("asset1", 3)
	investor.AddAssetPosition(investorAssetPosition)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)

	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(0)
	// investidor 1 quer vender 3 shares
	order := NewOrder("1", investor, asset1, 3, 6.0, "SELL")
	orderChan <- order

	// investidor 2 quer comprar 5 shares
	order2 := NewOrder("2", investor2, asset1, 5, 5.0, "BUY")
	orderChan <- order2

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	assert := assert.New(t)
	assert.Equal("OPEN", order.Status, "Order 1 should be closed")
	assert.Equal("OPEN", order2.Status, "Order 2 should be OPEN")
	assert.Equal(3, order.PendingShares, "Order 1 should have 3 PendingShares")
	assert.Equal(5, order2.PendingShares, "Order 2 should have 5 PendingShares")
}
