package todolist

import "github.com/rivo/tview"

type Component struct {
	app *tview.Application
}

func NewComponent(app *tview.Application) *Component {
	return &Component{
		app: app,
	}
}

func (c *Component) TodoList() *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false)
	list.SetBorder(true).SetTitle("TODOs")

	return list
}
