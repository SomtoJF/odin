package statusbar

import (
	"fmt"

	"github.com/SomtoJF/odin/core/state"
	"github.com/rivo/tview"
)

type Component struct {
	app *tview.Application
}

func NewComponent(app *tview.Application) *Component {
	return &Component{
		app: app,
	}
}

func (c *Component) StatusBar(currentMode state.AgentMode) *tview.TextView {
	statusText := fmt.Sprintf(
		"[Mode: [cyan]%s[white]] | Ctrl+A: Ask | Ctrl+E: Edit | Ctrl+P: Plan | Ctrl+Q: Quit",
		currentMode,
	)

	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText(statusText)
	statusBar.SetBorder(true).SetTitle("Status")

	return statusBar
}
