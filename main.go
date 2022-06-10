package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"typeracer-tui/text_view"
	"typeracer-tui/type_view"

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

	fileView := text_view.Init(app, fileContent, config.filePath)
	typeView := type_view.Init(app)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}

		if event.Key() == tcell.KeyRune {
			key := string(event.Rune())
			type_view.Refresh(type_view.EVENT_CHAR, key)
			text_view.Refresh(key)
		}

		if event.Key() == tcell.KeyEnter {
			type_view.Refresh(type_view.EVENT_NEW_LINE, "")
		}

		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			type_view.Refresh(type_view.EVENT_BACKSPACE, "")
		}

		return event
	})

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(fileView, 0, 1, false).
			AddItem(typeView, 0, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetBorderColor(tcell.ColorOlive).SetTitle("Stats"), 30, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
