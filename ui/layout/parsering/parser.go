package parser

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/ryo-del/devops-toolkit/internal/parser"
)

func NewParserTab() fyne.CanvasObject {

	infolabel := widget.NewLabel("Enter log file path:")
	entry := widget.NewEntry()
	entry.SetPlaceHolder("e.g., access.log")
	parlabel := widget.NewLabel("Parsed Log Entries:")
	output := widget.NewMultiLineEntry()
	output.Disable()

	parseBtn := widget.NewButton("Parse", func() {
		path := entry.Text
		if path == "" {
			output.SetText("Please enter a log file path.")
			return
		}
		entries, err := parser.ParseFile(path)
		if err != nil {
			output.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		var sb strings.Builder
		for _, e := range entries {
			sb.WriteString(fmt.Sprintf("%s %s %s %s %s\n", e.IP, e.Method, e.Path, e.StatusCode, e.ReqTime))
		}
		output.SetText(sb.String())
	})

	return container.NewVBox(
		infolabel,
		entry,
		parseBtn,
		parlabel,
		output,
	)

}
