package type_view

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var CursorName = "type_c"
var view *tview.TextView


func Init(app *tview.Application) *tview.TextView {
	newFileContent := Highlight(" ")
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
	fileView.SetTitle("Type...").SetBorder(true).SetBorderColor(tcell.ColorOliveDrab)
	view = fileView
	return fileView
}

func Refresh(event string, key string) {
	switch event {
	case EVENT_CHAR:
		text := view.GetText(false)
		view.Clear()
		fmt.Fprintf(view, "%s", AddNewKey(text, key))
	case EVENT_NEW_LINE:
		text := view.GetText(false)
		view.Clear()
		fmt.Fprintf(view, "%s", AddNewKey(text, "\n"))
	case EVENT_BACKSPACE:
		text := view.GetText(false)
		view.Clear()
		fmt.Fprintf(view, "%s", RemoveLastKey(text))
	}

}

func Highlight(character string) string {
	return fmt.Sprintf(`["%s"]%s[""]`, CursorName, character)
}

func RemoveEmptyHighlight(text string, cursorName string) string {
	emptyCursor := fmt.Sprintf(`["%s"] [""]`, cursorName)
	newText := strings.ReplaceAll(text, emptyCursor, "")
	newText = newText[:len(newText)-1]
	return newText
}

func AddNewKey(text string, key string) string {
	newTextWithoutHighlight := RemoveEmptyHighlight(text, CursorName)
	newText := newTextWithoutHighlight + key + Highlight(" ")
	return newText
}

func RemoveLastKey(text string) string {
	newText := RemoveEmptyHighlight(text, CursorName)
	if (len(newText) - 1) < 0 {
		return Highlight(" ")
	}
	newText = newText[:len(newText)-1]
	return newText + Highlight(" ")
}
