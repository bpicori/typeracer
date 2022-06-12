package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

type Config struct {
	filePath string
	language string
}

type Element struct {
	Text      string
	State     string
	Highlight bool
}

var States []Element
var textView *tview.TextView
var statsView *tview.TextView
var charIndex = 0
var timer = 60

func parseArguments() Config {
	filePath := flag.String("f", "", "file path")
	language := flag.String("l", "js", "language")
	flag.Parse()
	return Config{
		filePath: *filePath,
		language: *language,
	}
}

func main() {
	config := parseArguments()
	flag.Parse()

	if config.filePath == "" {
		fmt.Println("file path is empty")
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(config.filePath)
	if err != nil {
		fmt.Print("error reading file: ", err)
		os.Exit(1)
	}

	fileContent := string(b)

	app := tview.NewApplication()

	textView = Init(app, fileContent, config.filePath)
	statsView = tview.NewTextView().SetText(GetStats())

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}

		if event.Key() == tcell.KeyRune {
			key := string(event.Rune())
			Refresh(EVENT_CHAR, key)
		}

		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			Refresh(EVENT_BACKSPACE, "")
		}

		return event
	})

	go func() {
		for range time.Tick(1 * time.Second) {
			app.QueueUpdateDraw(func() {
				timer--
				statsView.Clear()
				fmt.Fprintf(statsView, "%s", GetStats())
			})
		}
	}()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(statsView, 0, 1, false).
		AddItem(textView, 0, 3, false)
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func GetStats() string {
	errors := 0
	counter := 0
	for i := 0; i < len(States); i++ {
		c := States[i]
		if c.State == STATE_WRONG && c.Text != " " {
			errors += 1
		}
		if c.State != STATE_UNDEFINED {
			counter += 1
		}
	}
	seconds := float64(60-timer) / 60.0
	wpm := (float64(counter-errors) / 5.0) / seconds
	return fmt.Sprintf("\n\ntimer: %d\nwpm: %f\nerrors: %d\n", timer, wpm, errors)
}

func Init(app *tview.Application, fileContent string, title string) *tview.TextView {
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

	fileView.SetTitle(title)
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
