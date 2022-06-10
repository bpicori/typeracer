package main

import (
	"fmt"
	"testing"
)

func TestRemoveHighlight(t *testing.T) {
	var tests = []struct {
		text       string
		cursorName string
		want       string
	}{
		{
			text:       " ",
			cursorName: "testCursor",
			want:       ` `,
		},
		{
			text:       `["testCursor"] [""]`,
			cursorName: "testCursor",
			want:       " ",
		},
		{
			text:       `test123["testCursor"]b[""]`,
			cursorName: "testCursor",
			want:       "test123b",
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s,%s", tt.text, tt.cursorName)
		t.Run(testName, func(t *testing.T) {
			ans := RemoveHighlight(tt.text, tt.cursorName)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

func TestAddKey(t *testing.T) {
	var tests = []struct {
		text string
		key  string
		want string
	}{
		{
			text: `bes["textCursor"]n[""]`,
			key:  "i",
			want: `besn["textCursor"]i[""]`,
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s,%s", tt.text, tt.key)
		t.Run(testName, func(t *testing.T) {
			ans := AddNewKey(tt.text, tt.key)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}

}
