package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	STATE_UNDEFINED = "undefined"
	STATE_CORRECT   = "correct"
	STATE_WRONG     = "wrong"
)

const (
	EVENT_NEW_LINE  = "new_line"
	EVENT_BACKSPACE = "backspace"
	EVENT_CHAR      = "character"
)

const (
	HIGHLIGHT_CURSOR = "cursor_c"
	GREEN_CURSOR     = "green_c"
	RED_CURSOR       = "red_c"
)

type Options struct {
	Timer int
	Words int
}

type Element struct {
	Text      string
	State     string
	Highlight bool
}

var States []Element
var app *tview.Application
var textView *tview.TextView
var statsView *tview.TextView
var charIndex = 0
var startTimer = false

func generateContent(nr int) string {
	rand.Seed(time.Now().UnixNano())
	var content string
	var words []string
	jsonFile, err := os.Open("language.json")
	if err != nil {
		os.Exit(1)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal([]byte(byteValue), &words)
	for i := 0; i < nr; i++ {
		index := rand.Intn(len(words))
		content += words[index] + " "
	}
	return content

}

func parseOptions() Options {
	words := flag.Int("w", 0, "number of words")
	timer := flag.Int("t", 0, "timer")

	flag.Parse()
	opts := Options{}
	if *words > 0 {
		opts.Words = *words
	} else {
		opts.Words = 100
	}
	if *timer > 0 {
		opts.Timer = *timer
	} else {
		opts.Timer = 60
	}
	return opts
}

func main() {
	opts := parseOptions()
	flag.Parse()

	content := generateContent(opts.Words)
	app = tview.NewApplication()

	textView = Init(app, content)
	statsView = tview.NewTextView().SetText(GetStats(opts.Timer))

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}

		if event.Key() == tcell.KeyRune {
			if !startTimer {
				StartTimer(opts)
				startTimer = true
			}
			key := string(event.Rune())
			Refresh(EVENT_CHAR, key)
		}

		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			Refresh(EVENT_BACKSPACE, "")
		}

		return event
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(statsView, 0, 1, false).
		AddItem(textView, 0, 3, false)
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func StartTimer(opts Options) {
	go func() {
		for range time.Tick(1 * time.Second) {
			app.QueueUpdateDraw(func() {
				opts.Timer--
				statsView.Clear()
				fmt.Fprintf(statsView, "%s", GetStats(opts.Timer))
			})
		}
	}()
}

func GetStats(timer int) string {
	if timer <= 0 {
		GameOver()
	}
	errors := 0
	counter := 0
	for i := 0; i < len(States); i++ {
		c := States[i]
		if c.State == STATE_WRONG && c.Text != " " {
			errors += 1
		}
		if c.State == STATE_CORRECT {
			counter += 1
		}
	}
	seconds := float64(60-timer) / 60.0
	wpm := (float64(counter) / 5.0) / seconds
	if math.IsNaN(wpm) {
		wpm = 0
	}
	return fmt.Sprintf("\n\ntimer: %d\nwpm: %d\nerrors: %d\n", timer, int(wpm), errors)
}

func Init(app *tview.Application, fileContent string) *tview.TextView {
	firstChar := fileContent[charIndex : charIndex+1]
	rest := fileContent[1:]
	newFileContent := Highlight(firstChar) + rest
	InitState(fileContent)
	fileView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetText(newFileContent).
		Highlight(HIGHLIGHT_CURSOR).
		SetTextAlign(tview.AlignCenter).
		SetToggleHighlights(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView = fileView

	return fileView
}

func InitState(fileContent string) {
	States = make([]Element, len(fileContent))
	for i := 0; i < len(fileContent); i++ {
		States[i].Text = fileContent[i : i+1]
		States[i].State = STATE_UNDEFINED
	}
}

func Render() string {
	text := ""
	for i := 0; i < len(States); i++ {
		c := States[i]
		if c.State == STATE_CORRECT {
			text += HighlightGreen(c.Text)
		} else if c.State == STATE_WRONG {
			text += HighlightRed(c.Text)
		} else if c.Highlight {
			text += Highlight(c.Text)
		} else {
			text += c.Text
		}
	}
	return text
}

func Refresh(event string, key string) {
	switch event {
	case EVENT_CHAR:
		textView.Clear()
		c := &States[charIndex]
		c.Highlight = false
		if c.Text == key {
			c.State = STATE_CORRECT
		} else {
			c.State = STATE_WRONG
		}
		charIndex += 1
		next := &States[charIndex]
		next.Highlight = true
		fmt.Fprintf(textView, "%s", Render())
	case EVENT_BACKSPACE:
		textView.Clear()
		current := &States[charIndex]
		current.Highlight = false
		previous := &States[charIndex-1]
		previous.State = STATE_UNDEFINED
		previous.Highlight = true
		charIndex -= 1
		fmt.Fprintf(textView, "%s", Render())
	}
	if charIndex+1 >= len(States) {
		GameOver()
	}
}

func GameOver() {
	modal := tview.NewModal().
		SetText("Game Over").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				app.Stop()
			}
		})
	app.SetRoot(modal, false)
	app.ForceDraw()
}

func Highlight(character string) string {
	return fmt.Sprintf(`["%s"]%s[""]`, HIGHLIGHT_CURSOR, character)
}

func HighlightGreen(character string) string {
	return fmt.Sprintf(`[green]%s[white]`, character)
}

func HighlightRed(character string) string {
	return fmt.Sprintf(`[red]%s[white]`, character)
}
