package messageview

import "github.com/rivo/tview"

type Component struct {
	app *tview.Application
}

func NewComponent(app *tview.Application) *Component {
	return &Component{
		app: app,
	}
}

func (c *Component) MessageView() *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			c.app.Draw()
		})
	textView.SetBorder(true).SetTitle("Messages")
	return textView
}
