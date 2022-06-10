package text_view

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var CursorName = "text_c"
var view *tview.TextView
var index = 0

func Init(app *tview.Application, fileContent string, title string) *tview.TextView {
	// highlight first char
	firstChar := fileContent[index : index+1]
	rest := fileContent[1:]
	newFileContent := highlight(firstChar) + rest
	fileView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetText(newFileContent).
		Highlight(CursorName).
		SetToggleHighlights(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	fileView.SetTitle(title).SetBorder(true).SetBorderColor(tcell.ColorOliveDrab)
	view = fileView

	return fileView
}

func Refresh(key string) {
	text := view.GetText(false)
	view.Clear()
	newText := RemoveHighlight(text)
	index += 1
	before := newText[:index]
	char := newText[index : index+1]
	rest := newText[index+1:]
	fmt.Fprintf(view, "%s", before+highlight(char)+rest)
}

func RemoveHighlight(text string) string {
	cursorPart1 := fmt.Sprintf(`["%s"]`, CursorName)
	cursorPart2 := `[""]`
	newText := strings.ReplaceAll(text, cursorPart1, "")
	newText = strings.ReplaceAll(newText, cursorPart2, "")
	newText = newText[:len(newText)-1]
	return newText
}

func highlight(character string) string {
	return fmt.Sprintf(`["%s"]%s[""]`, CursorName, character)
}
