package text_view

import (
	"fmt"

	"github.com/rivo/tview"
)

var (
	CursorName      = "cursor_c"
	GreenCursorName = "green_c"
	RedCursorName   = "red_c"
	index           = 0
)
var view *tview.TextView

type Element struct {
	Text      string
	State     string
	Highlight bool
}

var States []Element

const (
	EVENT_NEW_LINE  = "new_line"
	EVENT_BACKSPACE = "backspace"
	EVENT_CHAR      = "character"
)

const (
	STATE_UNDEFINED = "undefined"
	STATE_CORRECT   = "correct"
	STATE_WRONG     = "wrong"
)

func InitState(fileContent string) {
	States = make([]Element, len(fileContent))
	for i := 0; i < len(fileContent); i++ {
		States[i].Text = fileContent[i : i+1]
		States[i].State = STATE_UNDEFINED
	}
}

func Init(app *tview.Application, fileContent string, title string) *tview.TextView {
	// highlight first char
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

// func RemoveHighlight(text string) string {
// 	cursorPart1 := fmt.Sprintf(`["%s"]`, CursorName)
// 	cursorPart2 := `[""]`
// 	newText := strings.ReplaceAll(text, cursorPart1, "")
// 	newText = strings.ReplaceAll(newText, cursorPart2, "")
// 	newText = newText[:len(newText)-1]
// 	return newText
// }

// func RemoveColors(text string) string {
// 	redColor := `["red"]`
// 	greenColor := `["green"]`
// 	newText := strings.ReplaceAll(text, redColor, "")
// 	newText = strings.ReplaceAll(newText, greenColor, "")
// 	newText = newText[:len(newText)-1]
// 	return newText

// }

func Highlight(character string) string {
	return fmt.Sprintf(`["%s"]%s[""]`, CursorName, character)
}

func HighlightGreen(character string) string {
	return fmt.Sprintf(`[green]%s[white]`, character)
}

func HighlightRed(character string) string {
	return fmt.Sprintf(`[red]%s[white]`, character)
}
