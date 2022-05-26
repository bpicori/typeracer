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

	b, err := ioutil.ReadFile(config.filePath)
	if err != nil {
		fmt.Print("error reading file: ", err)
		os.Exit(1)
	}

	fileContent := string(b)
	cursor := `["cursor"] [""]`

	app := tview.NewApplication()

	fileView := getFileView(app, fileContent, cursor, config.filePath)
	typeView := getFileView(app, "", "", "Type..")

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}
		if event.Key() == tcell.KeyRune {
			key := string(event.Rune())
			text := fileView.GetText(false)
			newText := strings.ReplaceAll(text, cursor, "")
			newText = newText[:len(newText)-1]
			fileView.Clear()
			fmt.Fprintf(fileView, "%s", newText+key+cursor)
		}

		if event.Key() == tcell.KeyEnter {
			key := "\n"
			text := fileView.GetText(false)
			newText := strings.ReplaceAll(text, cursor, "")
			newText = newText[:len(newText)-1]
			fileView.Clear()
			fmt.Fprintf(fileView, "%s", newText+key+cursor)
		}

		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			text := fileView.GetText(false)
			newText := strings.ReplaceAll(text, cursor, "")
			if (len(newText) - 2) < 0 {
				fileView.Clear()
				fmt.Fprintf(fileView, "%s", cursor)
			} else {
				newText = newText[:len(newText)-2]
				fileView.Clear()
				fmt.Fprintf(fileView, "%s", newText+cursor)
			}
		}

		return event
	})

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(fileView, 0, 1, false).
			AddItem(typeView, 0, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetBorderColor(tcell.ColorOlive).SetTitle("Right (20 cols)"), 30, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func getFileView(app *tview.Application, fileContent string, cursor string, title string) *tview.TextView {
	fileView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetText(fileContent + cursor).
		Highlight("cursor").
		SetToggleHighlights(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	fileView.SetTitle(title).SetBorder(true).SetBorderColor(tcell.ColorOliveDrab)

	return fileView
}
