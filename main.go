// Demo code for the TextView primitive.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Config struct {
	filePath string
	language string
}

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

	_, err := ioutil.ReadFile(config.filePath)
	if err != nil {
		fmt.Print("error reading file: ", err)
		os.Exit(1)
	}

	// fileContent := string(b)
	cursor := `["cursor"] [""]`

	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetText(cursor).
		Highlight("cursor").
		SetToggleHighlights(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}
		if event.Key() == tcell.KeyRune {
			key := string(event.Rune())
			text := textView.GetText(false)
			newText := strings.ReplaceAll(text, cursor, "")
			newText = newText[:len(newText)-1]
			textView.Clear()
			fmt.Fprintf(textView, "%s", newText+key+cursor)
		}

		if event.Key() == tcell.KeyEnter {
			key := "\n"
			text := textView.GetText(false)
			newText := strings.ReplaceAll(text, cursor, "")
			newText = newText[:len(newText)-1]
			textView.Clear()
			fmt.Fprintf(textView, "%s", newText+key+cursor)
		}

		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			text := textView.GetText(false)
			newText := strings.ReplaceAll(text, cursor, "")
			newText = newText[:len(newText)-2]
			textView.Clear()
			fmt.Fprintf(textView, "%s", newText+cursor)
		}

		return event
	})

	textView.SetBorder(true)
	if err := app.SetRoot(textView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
