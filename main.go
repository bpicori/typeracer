package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Config struct {
	filePath string
	language string
}

var (
	fileCursorName = "fileCursor"
	textCursor     = "textCursor"
)

var (
	CursorName      = "cursor_c"
	GreenCursorName = "green_c"
	RedCursorName   = "red_c"
	index           = 0
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

type Element struct {
	Text      string
	State     string
	Highlight bool
}

var States []Element
var view *tview.TextView

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

	view = Init(app, fileContent, config.filePath)
	// typeView := type_view.Init(app)
	view.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		y += h / 2
		return x, y, w, h
	})

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

	if err := app.SetRoot(view, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func Init(app *tview.Application, fileContent string, title string) *tview.TextView {
	firstChar := fileContent[index : index+1]
	rest := fileContent[1:]
	newFileContent := Highlight(firstChar) + rest
	InitState(fileContent)
	fileView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetText(newFileContent).
		Highlight(CursorName).
		SetTextAlign(tview.AlignCenter).
		SetToggleHighlights(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	fileView.SetTitle(title)
	view = fileView

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
		view.Clear()
		c := &States[index]
		c.Highlight = false
		if c.Text == key {
			c.State = STATE_CORRECT
		} else {
			c.State = STATE_WRONG
		}
		index += 1
		next := &States[index]
		next.Highlight = true

		fmt.Fprintf(view, "%s", Render())
	case EVENT_BACKSPACE:
		view.Clear()
		current := &States[index]
		current.Highlight = false
		previous := &States[index-1]
		previous.State = STATE_UNDEFINED
		previous.Highlight = true
		index -= 1
		fmt.Fprintf(view, "%s", Render())
	}
}

func Highlight(character string) string {
	return fmt.Sprintf(`["%s"]%s[""]`, CursorName, character)
}

func HighlightGreen(character string) string {
	return fmt.Sprintf(`[green]%s[white]`, character)
}

func HighlightRed(character string) string {
	return fmt.Sprintf(`[red]%s[white]`, character)
}
