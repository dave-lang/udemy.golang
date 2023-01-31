package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var balancePlot *widgets.Plot
var ls *widgets.List
var lsLedger *widgets.List
var grid *ui.Grid

func buildUi(initialBalance int) {
	balancePlot = widgets.NewPlot()

	var labels []string
	for i := 1; i < 900; i++ {
		labels = append(labels, fmt.Sprint(i))
	}

	balancePlot.Marker = widgets.MarkerDot
	balancePlot.Title = "Balance (40 weeks)"
	balancePlot.Data = make([][]float64, 1)
	balancePlot.Data[0] = []float64{float64(initialBalance), float64(initialBalance)} // Needs 2 val to initialise
	balancePlot.SetRect(0, 0, 10, 10)
	balancePlot.DataLabels = labels
	balancePlot.Inner.Max.X = 30
	balancePlot.HorizontalScale = 1
	balancePlot.AxesColor = ui.ColorYellow
	balancePlot.LineColors[0] = ui.ColorGreen

	// Shows weekly balance
	ls = widgets.NewList()
	ls.Title = "Weekly balance (Up, down to scroll)"
	ls.Rows = []string{
		"Starting the week",
	}
	ls.Border = true

	// Shows bills and income
	lsLedger = widgets.NewList()
	lsLedger.Title = "Income and Bills (o, p to scroll)"
	lsLedger.Rows = []string{
		"Waiting for bills",
	}
	//lsLedger.Border = true

	grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(2.0/3, // 2/3 Height
			ui.NewCol(1.0/2, balancePlot), // 1/2 Width
			ui.NewCol(1.0/2, lsLedger),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0, ls),
		),
	)

	ui.Render(grid)
}

func updateBalance(balance float64) {
	balancePlot.Data[0] = append(balancePlot.Data[0], balance)

	// Limit to 40 weeks by slicing off first elemnt
	if len(balancePlot.Data[0]) > 40 {
		balancePlot.Data[0] = balancePlot.Data[0][1:]
	}

	ui.Render(balancePlot)
}

func updateLedger(str string) {
	lsLedger.Rows = append(lsLedger.Rows, str)
	//lsLedger.ScrollBottom()
	ui.Render(lsLedger)
}

func updateLog(str string) {
	ls.Rows = append(ls.Rows, str)
	ls.ScrollBottom()
	ui.Render(ls)
}

func handleInput() {
	previousKey := ""
	uiEvents := ui.PollEvents()
	//ticker := time.NewTicker(time.Second).C // Handle consistent updates

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "j", "<Down>":
				ls.ScrollDown()
				//lsLedger.ScrollDown()
			case "k", "<Up>":
				ls.ScrollUp()
				lsLedger.ScrollUp()
			case "<C-d>":
				ls.ScrollHalfPageDown()
			case "<C-u>":
				ls.ScrollHalfPageUp()
			case "<C-f>":
				ls.ScrollPageDown()
			case "<C-b>":
				ls.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					ls.ScrollTop()
				}
			case "o":
				lsLedger.ScrollUp()
			case "p":
				lsLedger.ScrollDown()
			case "<Home>":
				ls.ScrollTop()
			case "G", "<End>":
				ls.ScrollBottom()
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

			ui.Render(grid)
			//case <-ticker:
			//ui.Render(ls)
			//ui.Render(lsLedger)
			//ui.Render(balancePlot)
		}
	}
}
