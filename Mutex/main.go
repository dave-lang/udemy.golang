package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
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
var currentBalance int = 50

// Get paid
func getPaid(wg *sync.WaitGroup, mx *sync.RWMutex) {
	defer wg.Done()

	// Don't ask don't tell
	salary := rand.Intn(10)

	if salary > 0 {
		color.Green("You made $%d this week :)\n", salary)
		mx.Lock()
		currentBalance = currentBalance + salary
		mx.Unlock()
	}
}

// Try to pay all the bills
func payBills(week int, wg *sync.WaitGroup, mx *sync.RWMutex) {
	defer wg.Done()

	for _, bill := range bills {
		if (week % bill.EveryXWeeks) == 0 {
			color.Magenta("%s is due!\n", bill.Source)
			mx.Lock()
			newBalance := currentBalance
			newBalance = newBalance - bill.Amount
			currentBalance = newBalance
			mx.Unlock()
		}
	}
}

// Get the balance
func startWeek(week int, mx *sync.RWMutex) {
	mx.RLock()
	fmt.Printf("It's week %d, time to adult. You have $%d\n", week, currentBalance)
	mx.RUnlock()
}

// Find out when you go broke
func main() {
	rand.Seed(time.Now().UnixNano())

	week := 0

	defer func() {
		color.Red("You went broke on week %d\n", week)
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
	}
}
