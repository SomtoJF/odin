package inputfield

import "github.com/rivo/tview"

type Component struct {
	app *tview.Application
}

func NewComponent(app *tview.Application) *Component {
	return &Component{
		app: app,
	}
}

func (c *Component) InputField() *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)
	inputField.SetBorder(true).SetTitle("Input")

	return inputField
}
