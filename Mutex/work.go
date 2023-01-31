package main

import (
	"fmt"
	"math/rand"
	"sync"
	//"github.com/fatih/color"
)

// Get paid
func getPaid(wg *sync.WaitGroup, mx *sync.RWMutex) {
	defer wg.Done()

	// Don't ask don't tell
	salary := rand.Intn(10)

	if salary > 0 {
		//color.Green("You made $%d this week :)\n", salary)
		updateLedger(fmt.Sprint("You made $", salary, " this week :)"))

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
			//color.Magenta("%s is due!\n", bill.Source)
			updateLedger(fmt.Sprint(bill.Source, " is due!"))

			mx.Lock()
			newBalance := currentBalance
			newBalance = newBalance - bill.Amount
			currentBalance = newBalance
			mx.Unlock()
		}
	}
}
