package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
)

var balanceChart *linechart.LineChart
var balanceHistory []float64

var lsWeekly *text.Text
var lsBalance *text.Text

//var lsLedger *widgets.List

func createBalanceChart(initialBalance float64) {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
		linechart.YAxisCustomScale(0, 550),
		linechart.YAxisFormattedValues(linechart.ValueFormatterRound),
	)
	if err != nil {
		panic(err)
	}

	balanceHistory = []float64{initialBalance}

	if err := lc.Series("first", balanceHistory); err != nil {
		panic(err)
	}

	balanceChart = lc
}

func createWeeklyLedger() {
	text, err := text.New(text.RollContent())
	if err != nil {
		panic(err)
	}

	lsWeekly = text
}

func createBalanceLedger() {
	text, err := text.New(text.RollContent())
	if err != nil {
		panic(err)
	}

	lsBalance = text
}

func buildUi(initialBalance int) {
	// Balance chart
	createBalanceChart(float64(initialBalance))
	createWeeklyLedger()
	createBalanceLedger()
}

func runUi() {
	const redrawInterval = 250 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())

	t, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	builder := grid.New()

	builder.Add(
		grid.RowHeightPerc(50,
			grid.Widget(balanceChart, container.Border(linestyle.Light), container.BorderTitle("PRESS Q TO QUIT")),
		),
	)

	builder.Add(
		grid.ColWidthPerc(50,
			grid.Widget(lsBalance,
				container.Border(linestyle.Light), container.BorderTitle("Weekly balances"),
			),
		),
		grid.ColWidthPerc(50,
			grid.Widget(lsWeekly,
				container.Border(linestyle.Light), container.BorderTitle("Weekly credit/debits"),
			),
		),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		fmt.Errorf("builder.Build => %v", err)
	}

	c, err := container.New(t, gridOpts...)
	if err != nil {
		fmt.Errorf("container.New => %v", err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		panic(err)
	}
}

func updateBalance(balance float64) {
	balanceHistory = append(balanceHistory, balance)

	if err := balanceChart.Series("first", balanceHistory,
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(33))),
		linechart.SeriesXLabels(map[int]string{
			0: "0",
		}),
	); err != nil {
		panic(err)
	}

	// // Limit to 40 weeks by slicing off first elemnt
	// if len(balancePlot.Data[0]) > 40 {
	// 	balancePlot.Data[0] = balancePlot.Data[0][1:]
	// }

	// ui.Render(balancePlot)
}

func updateLedger(str string) {
	if err := lsBalance.Write(fmt.Sprintf("%s\n", str)); err != nil {
		panic(err)
	}
}

func updateWeekly(str string) {
	if err := lsWeekly.Write(fmt.Sprintf("%s\n", str)); err != nil {
		panic(err)
	}
}
