package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//Building web-enabled Terminal UI apps in Go

const refreshInterval = 10 * time.Second
const url = "https://api.chucknorris.io/jokes/random?category=science"

// Global Variables
var (
	app      *tview.Application
	textView *tview.TextView
)

type Payload struct {
	Value string `json:"value"`
}

func getAndDrawJoke() {
	// fetch chuck norris joke from the web
	result, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	payloadBytes, err := io.ReadAll(result.Body)

	if err != nil {
		panic(err)
	}

	payload := &Payload{}
	err = json.Unmarshal(payloadBytes, payload)

	if err != nil {
		panic(err)
	}

	// update our UI with the joke

	textView.Clear()
	fmt.Fprintln(textView, payload.Value)
	timeStr := fmt.Sprintf("\n\n[gray]%s", time.Now().Format(time.RFC1123))
	fmt.Fprintln(textView, timeStr)
}

// Refresh the Joke

func refreshJoke() {

	tick := time.NewTicker(refreshInterval)

	for {
		select {
		case <-tick.C:
			getAndDrawJoke()
			app.Draw()
		}
	}
}

func renderFooter() *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Built with [red]Love[white] by Jubilio Mausse ([gray]@jubiliomausse[white])")
}

func renderHeader() *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText(`
		[yellow]Chuck Norris Jokes
		[white] We Love Jokes!
		[red] Science Jokes Only!
		`)
}

// Main App
func main() {
	app = tview.NewApplication()

	textView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorLime).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorderPadding(1, 0, 0, 0)

	getAndDrawJoke()

	grid := tview.NewGrid().
		SetRows(10, 0, 3).
		SetColumns(0, 0).
		AddItem(renderHeader(), 0, 0, 1, 2, 0, 0, false).
		AddItem(renderFooter(), 2, 0, 1, 2, 0, 0, false)

	grid.AddItem(textView, 1, 0, 1, 2, 0, 0, false)

	go refreshJoke()

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
