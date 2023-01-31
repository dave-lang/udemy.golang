package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Bill struct {
	Source      string
	EveryXWeeks int
	Amount      int
}

// Bills to be paid
var bills []Bill = []Bill{
	{Source: "Rent", EveryXWeeks: 2, Amount: 20},
}

// Starting balance
var currentBalance int = 500

// Get the balance
func startWeek(week int, mx *sync.RWMutex) {
	mx.RLock()
	updateLedger(fmt.Sprint("It's week ", week, ", you have $", currentBalance))
	updateBalance(float64(currentBalance))
	mx.RUnlock()
}

// Find out when you go broke
func main() {
	rand.Seed(time.Now().UnixNano())

	buildUi(currentBalance)
	go simulate()
	runUi()
}

// Actually does money
func simulate() {
	time.Sleep(time.Duration(1 * time.Second))

	week := 0

	// When did we fail at life?
	defer func() {
		//color.Red("You went broke on week %d\n", week)
		updateLedger(fmt.Sprint("You went broke on week", week))
	}()

	var wg sync.WaitGroup
	var mx sync.RWMutex

	for currentBalance > 0 {
		week++

		startWeek(week, &mx)

		wg.Add(2) // We do 3 things every week
		go getPaid(&wg, &mx)
		go payBills(week, &wg, &mx)
		wg.Wait()

		time.Sleep(time.Duration(1 * time.Second))
	}
}
