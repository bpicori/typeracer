package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"typeracer-tui/text_view"

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
	// typeView := type_view.Init(app)
	fileView.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		y += h / 2
		return x, y, w, h
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}

		if event.Key() == tcell.KeyRune {
			key := string(event.Rune())
			text_view.Refresh(text_view.EVENT_CHAR, key)
		}

		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			text_view.Refresh(text_view.EVENT_BACKSPACE, "")
		}

		return event
	})

	if err := app.SetRoot(fileView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
