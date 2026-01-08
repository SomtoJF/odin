package components

import (
	"github.com/SomtoJF/odin/ui/components"
	"github.com/rivo/tview"
)

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) Run() {
	box := tview.NewBox().SetBorder(true).SetTitle(components.IntroASCII())
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
